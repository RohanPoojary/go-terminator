package terminator

import (
	"context"
	"os"
	"time"
)

type TerminationStatus string

const (
	SUCCESS TerminationStatus = "SUCCESS"
	FAILED  TerminationStatus = "FAILED"
)

type TerminationResultData struct {
	Name   string
	Error  error
	Status TerminationStatus
}

type TerminationResult struct {
	Signal               os.Signal
	FailedOrTimeoutCount int // Number of Failed or Timeout Instances
	Result               []TerminationResultData
}

type CloseFunc func(context.Context) error

type Terminator interface {
	Add(name string, close CloseFunc)
	AddWithTimeout(name string, close CloseFunc, timeout time.Duration)
	SetCallback(callback func(TerminationResult))
	Wait(timeout time.Duration) bool
}
