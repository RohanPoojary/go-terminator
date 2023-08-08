package terminator

import (
	"context"
	"os"
	"time"
)

// TerminationStatus represents the status of termination (success or failure).
type TerminationStatus string

const (

	// SUCCESS indicates that the resource was closed successfully.
	SUCCESS TerminationStatus = "SUCCESS"

	// FAILED indicates that the resource failed to close.
	FAILED TerminationStatus = "FAILED"
)

// TerminationResultData holds information about the result of terminating a resource.
type TerminationResultData struct {

	// Name of the terminated resource
	Name string

	// Error that occurred during termination, if any
	Error error

	// Termination status of the process
	Status TerminationStatus
}

// TerminationResult contains the overall result of the termination process.
type TerminationResult struct {

	// Termination signal received
	Signal os.Signal

	// Number of resources that failed or timed out
	FailedOrTimeoutCount int

	// Result data for each terminated resource
	Result []TerminationResultData
}

// CloseFunc defines the function signature for closing a resource.
type CloseFunc func(context.Context) error

// Terminator is the interface that provides methods for managing resource termination.
type Terminator interface {

	// Add registers a resource to be closed without a timeout.
	Add(name string, close CloseFunc)

	// AddWithTimeout registers a resource to be closed with a specified timeout.
	AddWithTimeout(name string, close CloseFunc, timeout time.Duration)

	// SetCallback sets the callback function to be executed after all resources are closed.
	SetCallback(callback func(TerminationResult))

	// Wait waits for the termination process to complete within the specified timeout duration.
	Wait(timeout time.Duration) bool
}
