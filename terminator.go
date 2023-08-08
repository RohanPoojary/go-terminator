// Package terminator provides a utility to gracefully terminate processes by closing registered resources when a termination signal is received.

package terminator

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// payload represents a resource that needs to be closed gracefully.
type payload struct {
	Name    string
	Timeout time.Duration
	Close   func(context.Context) error
}

type terminator struct {
	closersStack  []payload
	signalChan    chan os.Signal
	completedChan chan bool
	callbackFunc  func(TerminationResult)
}

// NewTerminator creates a new instance of the terminator.
func NewTerminator(closeSignals []os.Signal) Terminator {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, closeSignals...)

	term := &terminator{
		signalChan:    sigc,
		completedChan: make(chan bool, 1),
	}

	go term.startMonitor()

	return term
}

// Add registers a resource with the terminator to be closed without any timeout.
func (t *terminator) Add(name string, close CloseFunc) {
	t.AddWithTimeout(name, close, 0)
}

// AddWithTimeout registers a resource with the terminator to be closed with a specified timeout.
func (t *terminator) AddWithTimeout(name string, close CloseFunc, timeout time.Duration) {
	t.closersStack = append(t.closersStack, payload{Name: name, Close: close, Timeout: timeout})
}

// SetCallback sets the callback function to be executed after all resources are closed.
func (t *terminator) SetCallback(fn func(TerminationResult)) {
	t.callbackFunc = fn
}

// Wait waits for the termination process to complete with a specified timeout duration.
func (t *terminator) Wait(timeout time.Duration) bool {
	select {
	case <-t.completedChan:
		return true
	case <-time.After(timeout):
		return false
	}
}

// closeStack performs the actual closing of a single resource in a separate goroutine.
func (t *terminator) closeStack(closer *payload) <-chan TerminationResultData {
	result := make(chan TerminationResultData, 1)

	ctx := context.Background()

	go func() {
		name := closer.Name
		// Apply timeout to the resource's closing if specified.
		if closer.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, closer.Timeout)
			defer cancel()
		}

		var status TerminationStatus
		var err error

		errChan := make(chan error, 1)
		select {
		case <-ctx.Done():
			err = ctx.Err()
			// If context has no error, then run close again.
			if err == nil {
				err = closer.Close(ctx)
			}
		case errChan <- closer.Close(ctx):
			err = <-errChan
		}

		if err == nil {
			status = SUCCESS
		} else {
			status = FAILED
		}

		result <- TerminationResultData{
			Name:   name,
			Status: status,
			Error:  err,
		}

	}()

	return result
}

// closeAll closes all the registered resources and collects the termination result data.
func (t *terminator) closeAll(ctx context.Context, result *TerminationResult) {

	var stackIndex int

	for stackIndex = len(t.closersStack) - 1; stackIndex >= 0; stackIndex-- {

		termData := <-t.closeStack(&t.closersStack[stackIndex])

		if termData.Error != nil {
			result.FailedOrTimeoutCount++
		}

		result.Result = append(result.Result, termData)
	}

}

// unsubscribe stops listening to termination signals.
func (t *terminator) unsubscribe() {
	signal.Stop(t.signalChan)
}

// startMonitor starts monitoring for termination signals and initiates the termination process.
func (t *terminator) startMonitor() {

	s := <-t.signalChan

	// Initializing Result
	result := TerminationResult{
		Signal: s,
		Result: make([]TerminationResultData, 0, len(t.closersStack)),
	}

	ctx := context.Background()

	t.closeAll(ctx, &result)

	if t.callbackFunc != nil {
		t.callbackFunc(result)
	}

	t.unsubscribe()
	close(t.completedChan)
}
