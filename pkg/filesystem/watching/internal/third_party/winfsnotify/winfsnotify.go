// Windows filesystem monitoring implementation based on
// golang.org/x/exp/winfsnotify
// (specifically
// https://github.com/golang/exp/tree/c84be7c6d1cd7b6a43fd7101daaf2dc35ded445f/winfsnotify),
// but modified to remove import path enforcement, increase
// ReadDirectoryChangesW buffer size, support recursive watching, use more
// idiomatic filesystem path joins, and remove test logging.
//
// The original code license:
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// The original license header inside the code itself:
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

// Package winfsnotify allows the user to receive
// file system event notifications on Windows.
package winfsnotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

// Event is the type of the notification messages
// received on the watcher's Event channel.
type Event struct {
	Mask   uint32 // Mask of events
	Cookie uint32 // Unique cookie associating related events (for rename)
	Name   string // File name (optional)
}

const (
	opAddWatch = iota
	opRemoveWatch
)

const (
	provisional uint64 = 1 << (32 + iota)
)

type input struct {
	op    int
	path  string
	flags uint32
	reply chan error
}

type inode struct {
	handle syscall.Handle
	volume uint32
	index  uint64
}

type watch struct {
	ov     syscall.Overlapped
	ino    *inode            // i-number
	path   string            // Directory path
	mask   uint64            // Directory itself is being watched with these notify flags
	names  map[string]uint64 // Map of names being watched and their notify flags
	rename string            // Remembers the old name while renaming a file
	buf    [16384]byte
}

type indexMap map[uint64]*watch
type watchMap map[uint32]indexMap

// A Watcher waits for and receives event notifications
// for a specific set of files and directories.
type Watcher struct {
	port     syscall.Handle // Handle to completion port
	watches  watchMap       // Map of watches (key: i-number)
	input    chan *input    // Inputs to the reader are sent on this channel
	Event    chan *Event    // Events are returned on this channel
	Error    chan error     // Errors are sent on this channel
	isClosed bool           // Set to true when Close() is first called
	quit     chan chan<- error
	cookie   uint32
}

// NewWatcher creates and returns a Watcher.
func NewWatcher() (*Watcher, error) {
	port, e := syscall.CreateIoCompletionPort(syscall.InvalidHandle, 0, 0, 0)
	if e != nil {
		return nil, os.NewSyscallError("CreateIoCompletionPort", e)
	}
	w := &Watcher{
		port:    port,
		watches: make(watchMap),
		input:   make(chan *input, 1),
		Event:   make(chan *Event, 50),
		Error:   make(chan error, 1),
		quit:    make(chan chan<- error, 1),
	}
	go w.readEvents()
	return w, nil
}

// Close closes a Watcher.
// It sends a message to the reader goroutine to quit and removes all watches
// associated with the watcher.
func (w *Watcher) Close() error {
	if w.isClosed {
		return nil
	}
	w.isClosed = true

	// Send "quit" message to the reader goroutine
	ch := make(chan error)
	w.quit <- ch
	if err := w.wakeupReader(); err != nil {
		return err
	}
	return <-ch
}

// AddWatch adds path to the watched file set.
func (w *Watcher) AddWatch(path string, flags uint32) error {
	if w.isClosed {
		return errors.New("watcher already closed")
	}
	in := &input{
		op:    opAddWatch,
		path:  filepath.Clean(path),
		flags: flags,
		reply: make(chan error),
	}
	w.input <- in
	if err := w.wakeupReader(); err != nil {
		return err
	}
	return <-in.reply
}

// Watch adds path to the watched file set, watching all events.
func (w *Watcher) Watch(path string) error {
	return w.AddWatch(path, FS_ALL_EVENTS)
}

// RemoveWatch removes path from the watched file set.
func (w *Watcher) RemoveWatch(path string) error {
	in := &input{
		op:    opRemoveWatch,
		path:  filepath.Clean(path),
		reply: make(chan error),
	}
	w.input <- in
	if err := w.wakeupReader(); err != nil {
		return err
	}
	return <-in.reply
}

func (w *Watcher) wakeupReader() error {
	e := syscall.PostQueuedCompletionStatus(w.port, 0, 0, nil)
	if e != nil {
		return os.NewSyscallError("PostQueuedCompletionStatus", e)
	}
	return nil
}

func getDir(pathname string) (dir string, err error) {
	pathnamep, e := syscall.UTF16PtrFromString(pathname)
	if e != nil {
		return "", e
	}
	attr, e := syscall.GetFileAttributes(pathnamep)
	if e != nil {
		return "", os.NewSyscallError("GetFileAttributes", e)
	}
	if attr&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 {
		dir = pathname
	} else {
		dir, _ = filepath.Split(pathname)
		dir = filepath.Clean(dir)
	}
	return
}

func getIno(path string) (ino *inode, err error) {
	pathp, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return nil, e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_LIST_DIRECTORY,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil, syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OVERLAPPED, 0)
	if e != nil {
		return nil, os.NewSyscallError("CreateFile", e)
	}
	var fi syscall.ByHandleFileInformation
	if e = syscall.GetFileInformationByHandle(h, &fi); e != nil {
		syscall.CloseHandle(h)
		return nil, os.NewSyscallError("GetFileInformationByHandle", e)
	}
	ino = &inode{
		handle: h,
		volume: fi.VolumeSerialNumber,
		index:  uint64(fi.FileIndexHigh)<<32 | uint64(fi.FileIndexLow),
	}
	return ino, nil
}

// Must run within the I/O thread.
func (m watchMap) get(ino *inode) *watch {
	if i := m[ino.volume]; i != nil {
		return i[ino.index]
	}
	return nil
}

// Must run within the I/O thread.
func (m watchMap) set(ino *inode, watch *watch) {
	i := m[ino.volume]
	if i == nil {
		i = make(indexMap)
		m[ino.volume] = i
	}
	i[ino.index] = watch
}

// Must run within the I/O thread.
func (w *Watcher) addWatch(pathname string, flags uint64) error {
	dir, err := getDir(pathname)
	if err != nil {
		return err
	}
	if flags&FS_ONLYDIR != 0 && pathname != dir {
		return nil
	}
	ino, err := getIno(dir)
	if err != nil {
		return err
	}
	watchEntry := w.watches.get(ino)
	if watchEntry == nil {
		if _, e := syscall.CreateIoCompletionPort(ino.handle, w.port, 0, 0); e != nil {
			syscall.CloseHandle(ino.handle)
			return os.NewSyscallError("CreateIoCompletionPort", e)
		}
		watchEntry = &watch{
			ino:   ino,
			path:  dir,
			names: make(map[string]uint64),
		}
		w.watches.set(ino, watchEntry)
		flags |= provisional
	} else {
		syscall.CloseHandle(ino.handle)
	}
	if pathname == dir {
		watchEntry.mask |= flags
	} else {
		watchEntry.names[filepath.Base(pathname)] |= flags
	}
	if err = w.startRead(watchEntry); err != nil {
		return err
	}
	if pathname == dir {
		watchEntry.mask &= ^provisional
	} else {
		watchEntry.names[filepath.Base(pathname)] &= ^provisional
	}
	return nil
}

// Must run within the I/O thread.
func (w *Watcher) removeWatch(pathname string) error {
	dir, err := getDir(pathname)
	if err != nil {
		return err
	}
	ino, err := getIno(dir)
	if err != nil {
		return err
	}
	watch := w.watches.get(ino)
	if watch == nil {
		return fmt.Errorf("can't remove non-existent watch for: %s", pathname)
	}
	if pathname == dir {
		w.sendEvent(watch.path, watch.mask&FS_IGNORED)
		watch.mask = 0
	} else {
		name := filepath.Base(pathname)
		w.sendEvent(filepath.Join(watch.path, name), watch.names[name]&FS_IGNORED)
		delete(watch.names, name)
	}
	return w.startRead(watch)
}

// Must run within the I/O thread.
func (w *Watcher) deleteWatch(watch *watch) {
	for name, mask := range watch.names {
		if mask&provisional == 0 {
			w.sendEvent(filepath.Join(watch.path, name), mask&FS_IGNORED)
		}
		delete(watch.names, name)
	}
	if watch.mask != 0 {
		if watch.mask&provisional == 0 {
			w.sendEvent(watch.path, watch.mask&FS_IGNORED)
		}
		watch.mask = 0
	}
}

// Must run within the I/O thread.
func (w *Watcher) startRead(watch *watch) error {
	if e := syscall.CancelIo(watch.ino.handle); e != nil {
		w.Error <- os.NewSyscallError("CancelIo", e)
		w.deleteWatch(watch)
	}
	mask := toWindowsFlags(watch.mask)
	for _, m := range watch.names {
		mask |= toWindowsFlags(m)
	}
	if mask == 0 {
		if e := syscall.CloseHandle(watch.ino.handle); e != nil {
			w.Error <- os.NewSyscallError("CloseHandle", e)
		}
		delete(w.watches[watch.ino.volume], watch.ino.index)
		return nil
	}
	e := syscall.ReadDirectoryChanges(watch.ino.handle, &watch.buf[0],
		uint32(unsafe.Sizeof(watch.buf)), true, mask, nil, &watch.ov, 0)
	if e != nil {
		err := os.NewSyscallError("ReadDirectoryChanges", e)
		if e == syscall.ERROR_ACCESS_DENIED && watch.mask&provisional == 0 {
			// Watched directory was probably removed
			if w.sendEvent(watch.path, watch.mask&FS_DELETE_SELF) {
				if watch.mask&FS_ONESHOT != 0 {
					watch.mask = 0
				}
			}
			err = nil
		}
		w.deleteWatch(watch)
		w.startRead(watch)
		return err
	}
	return nil
}

// readEvents reads from the I/O completion port, converts the
// received events into Event objects and sends them via the Event channel.
// Entry point to the I/O thread.
func (w *Watcher) readEvents() {
	var (
		n, key uint32
		ov     *syscall.Overlapped
	)
	runtime.LockOSThread()

	for {
		e := syscall.GetQueuedCompletionStatus(w.port, &n, &key, &ov, syscall.INFINITE)
		watch := (*watch)(unsafe.Pointer(ov))

		if watch == nil {
			select {
			case ch := <-w.quit:
				for _, index := range w.watches {
					for _, watch := range index {
						w.deleteWatch(watch)
						w.startRead(watch)
					}
				}
				var err error
				if e := syscall.CloseHandle(w.port); e != nil {
					err = os.NewSyscallError("CloseHandle", e)
				}
				close(w.Event)
				close(w.Error)
				ch <- err
				return
			case in := <-w.input:
				switch in.op {
				case opAddWatch:
					in.reply <- w.addWatch(in.path, uint64(in.flags))
				case opRemoveWatch:
					in.reply <- w.removeWatch(in.path)
				}
			default:
			}
			continue
		}

		switch e {
		case syscall.ERROR_ACCESS_DENIED:
			// Watched directory was probably removed
			w.sendEvent(watch.path, watch.mask&FS_DELETE_SELF)
			w.deleteWatch(watch)
			w.startRead(watch)
			continue
		case syscall.ERROR_OPERATION_ABORTED:
			// CancelIo was called on this handle
			continue
		default:
			w.Error <- os.NewSyscallError("GetQueuedCompletionPort", e)
			continue
		case nil:
		}

		var offset uint32
		for {
			if n == 0 {
				w.Event <- &Event{Mask: FS_Q_OVERFLOW}
				w.Error <- errors.New("short read in readEvents()")
				break
			}

			// Point "raw" to the event in the buffer
			raw := (*syscall.FileNotifyInformation)(unsafe.Pointer(&watch.buf[offset]))
			buf := (*[syscall.MAX_PATH]uint16)(unsafe.Pointer(&raw.FileName))
			name := syscall.UTF16ToString(buf[:raw.FileNameLength/2])
			fullname := filepath.Join(watch.path, name)

			var mask uint64
			switch raw.Action {
			case syscall.FILE_ACTION_REMOVED:
				mask = FS_DELETE_SELF
			case syscall.FILE_ACTION_MODIFIED:
				mask = FS_MODIFY
			case syscall.FILE_ACTION_RENAMED_OLD_NAME:
				watch.rename = name
			case syscall.FILE_ACTION_RENAMED_NEW_NAME:
				if watch.names[watch.rename] != 0 {
					watch.names[name] |= watch.names[watch.rename]
					delete(watch.names, watch.rename)
					mask = FS_MOVE_SELF
				}
			}

			sendNameEvent := func() {
				if w.sendEvent(fullname, watch.names[name]&mask) {
					if watch.names[name]&FS_ONESHOT != 0 {
						delete(watch.names, name)
					}
				}
			}
			if raw.Action != syscall.FILE_ACTION_RENAMED_NEW_NAME {
				sendNameEvent()
			}
			if raw.Action == syscall.FILE_ACTION_REMOVED {
				w.sendEvent(fullname, watch.names[name]&FS_IGNORED)
				delete(watch.names, name)
			}
			if w.sendEvent(fullname, watch.mask&toFSnotifyFlags(raw.Action)) {
				if watch.mask&FS_ONESHOT != 0 {
					watch.mask = 0
				}
			}
			if raw.Action == syscall.FILE_ACTION_RENAMED_NEW_NAME {
				fullname = filepath.Join(watch.path, watch.rename)
				sendNameEvent()
			}

			// Move to the next event in the buffer
			if raw.NextEntryOffset == 0 {
				break
			}
			offset += raw.NextEntryOffset
		}

		if err := w.startRead(watch); err != nil {
			w.Error <- err
		}
	}
}

func (w *Watcher) sendEvent(name string, mask uint64) bool {
	if mask == 0 {
		return false
	}
	event := &Event{Mask: uint32(mask), Name: name}
	if mask&FS_MOVE != 0 {
		if mask&FS_MOVED_FROM != 0 {
			w.cookie++
		}
		event.Cookie = w.cookie
	}
	select {
	case ch := <-w.quit:
		w.quit <- ch
	case w.Event <- event:
	}
	return true
}

// String formats the event e in the form
// "filename: 0xEventMask = FS_ACCESS|FS_ATTRIB_|..."
func (e *Event) String() string {
	var events string
	m := e.Mask
	for _, b := range eventBits {
		if m&b.Value != 0 {
			m &^= b.Value
			events += "|" + b.Name
		}
	}
	if m != 0 {
		events += fmt.Sprintf("|%#x", m)
	}
	if len(events) > 0 {
		events = " == " + events[1:]
	}
	return fmt.Sprintf("%q: %#x%s", e.Name, e.Mask, events)
}

func toWindowsFlags(mask uint64) uint32 {
	var m uint32
	if mask&FS_ACCESS != 0 {
		m |= syscall.FILE_NOTIFY_CHANGE_LAST_ACCESS
	}
	if mask&FS_MODIFY != 0 {
		m |= syscall.FILE_NOTIFY_CHANGE_LAST_WRITE
	}
	if mask&FS_ATTRIB != 0 {
		m |= syscall.FILE_NOTIFY_CHANGE_ATTRIBUTES
	}
	if mask&(FS_MOVE|FS_CREATE|FS_DELETE) != 0 {
		m |= syscall.FILE_NOTIFY_CHANGE_FILE_NAME | syscall.FILE_NOTIFY_CHANGE_DIR_NAME
	}
	return m
}

func toFSnotifyFlags(action uint32) uint64 {
	switch action {
	case syscall.FILE_ACTION_ADDED:
		return FS_CREATE
	case syscall.FILE_ACTION_REMOVED:
		return FS_DELETE
	case syscall.FILE_ACTION_MODIFIED:
		return FS_MODIFY
	case syscall.FILE_ACTION_RENAMED_OLD_NAME:
		return FS_MOVED_FROM
	case syscall.FILE_ACTION_RENAMED_NEW_NAME:
		return FS_MOVED_TO
	}
	return 0
}

const (
	// Options for AddWatch
	FS_ONESHOT = 0x80000000
	FS_ONLYDIR = 0x1000000

	// Events
	FS_ACCESS      = 0x1
	FS_ALL_EVENTS  = 0xfff
	FS_ATTRIB      = 0x4
	FS_CLOSE       = 0x18
	FS_CREATE      = 0x100
	FS_DELETE      = 0x200
	FS_DELETE_SELF = 0x400
	FS_MODIFY      = 0x2
	FS_MOVE        = 0xc0
	FS_MOVED_FROM  = 0x40
	FS_MOVED_TO    = 0x80
	FS_MOVE_SELF   = 0x800

	// Special events
	FS_IGNORED    = 0x8000
	FS_Q_OVERFLOW = 0x4000
)

var eventBits = []struct {
	Value uint32
	Name  string
}{
	{FS_ACCESS, "FS_ACCESS"},
	{FS_ATTRIB, "FS_ATTRIB"},
	{FS_CREATE, "FS_CREATE"},
	{FS_DELETE, "FS_DELETE"},
	{FS_DELETE_SELF, "FS_DELETE_SELF"},
	{FS_MODIFY, "FS_MODIFY"},
	{FS_MOVED_FROM, "FS_MOVED_FROM"},
	{FS_MOVED_TO, "FS_MOVED_TO"},
	{FS_MOVE_SELF, "FS_MOVE_SELF"},
	{FS_IGNORED, "FS_IGNORED"},
	{FS_Q_OVERFLOW, "FS_Q_OVERFLOW"},
}
