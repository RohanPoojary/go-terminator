package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/RohanPoojary/go-terminator"
)

func main() {
	term := terminator.NewTerminator([]os.Signal{os.Interrupt})
	defer term.Wait(5 * time.Second)

	fmt.Println("Hello, playground")

	// interruptCount := 0
	// term.Callback(func(result terminator.TerminationResult) {
	// 	interruptCount++
	// 	fmt.Println("callback called tried:", interruptCount)
	// 	fmt.Println(result)
	// })

	term.AddWithTimeout("app1", func(ctx context.Context) error {
		fmt.Println("closing app 1")
		time.Sleep(1 * time.Second)
		return nil
	}, time.Millisecond*500)

	term.Add("app2", func(ctx context.Context) error {
		fmt.Println("closing app 2")
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	term.Add("app3", func(ctx context.Context) error {
		fmt.Println("closing app 3")
		time.Sleep(200 * time.Millisecond)
		return errors.New("Failed App3")
	})

	closeLoop := make(chan bool, 1)
	term.Add("Loop Closer", func(ctx context.Context) error {
		closeLoop <- true
		return nil
	})

	for {
		select {
		case <-closeLoop:
			fmt.Println("closing loop")
			goto breakLoop
		case <-time.After(1 * time.Second):
			fmt.Println("sleeping")
		}
	}

breakLoop:
}
