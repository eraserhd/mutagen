package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mutagen-io/mutagen/cmd"

	"github.com/mutagen-io/mutagen/pkg/filesystem/watching"
)

const (
	// eventsBufferSize is the buffer size to use for the events channel.
	eventsBufferSize = 50
)

func main() {
	// Verify that the platform supports recursive watching.
	if !watching.RecursiveWatchingSupported {
		cmd.Fatal(errors.New("recursive watching not supported"))
	}

	// Parse arguments.
	if len(os.Args) != 2 {
		cmd.Fatal(errors.New("invalid number of arguments"))
	}
	watchRoot := os.Args[1]

	// Print status.
	fmt.Println("Watching", watchRoot)

	// Create an events channel.
	events := make(chan string, eventsBufferSize)

	// Perform watching in the background and track any errors. Close the events
	// channel once watching terminates.
	watchErrors := make(chan error, 1)
	go func() {
		watchErrors <- watching.WatchRecursive(
			context.Background(),
			watchRoot,
			events,
		)
		close(events)
	}()

	// Print events until watching has terminated.
	for path := range events {
		fmt.Printf("Event: \"%s\"\n", path)
	}

	// Wait for the watch error.
	cmd.Fatal(fmt.Errorf("watching failed: %w", <-watchErrors))
}
