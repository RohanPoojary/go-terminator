package terminator

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestBasicShutdown(t *testing.T) {
	term := NewTerminator([]os.Signal{os.Interrupt})

	result := []string{}
	term.Add("app1", func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		result = append(result, "app1")
		return nil
	})

	termInternal := term.(*terminator)
	termInternal.signalChan <- os.Interrupt

	ok := term.Wait(5 * time.Second)
	if !ok {
		t.Error("Wait timed out")
		return
	}

	if result[0] != "app1" {
		t.Error("app1 not closed")
	}
}

func TestWaitTimeout(t *testing.T) {
	term := NewTerminator([]os.Signal{os.Interrupt})

	result := []string{}
	term.Add("app1", func(ctx context.Context) error {
		time.Sleep(5 * time.Second)
		result = append(result, "app1")
		return nil
	})

	termInternal := term.(*terminator)
	termInternal.signalChan <- os.Interrupt

	ok := term.Wait(1 * time.Second)
	if ok {
		t.Error("Wait should have timed out")
		return
	}
}

func TestExecutionOrder(t *testing.T) {
	term := NewTerminator([]os.Signal{os.Interrupt})

	result := []string{}

	term.Add("app1", func(ctx context.Context) error {
		result = append(result, "app1")
		return nil
	})

	term.Add("app2", func(ctx context.Context) error {
		result = append(result, "app2")
		return nil
	})

	term.Add("app3", func(ctx context.Context) error {
		result = append(result, "app3")
		return nil
	})

	termInternal := term.(*terminator)
	termInternal.signalChan <- os.Interrupt

	ok := term.Wait(1 * time.Second)
	if !ok {
		t.Error("Wait shouldn't time out")
		return
	}

	if len(result) != 3 {
		t.Error("All apps should have been closed")
		return
	}

	if result[2] != "app1" || result[1] != "app2" || result[0] != "app3" {
		t.Error("Execution order not maintained")
		return
	}
}

func TestCallback(t *testing.T) {
	term := NewTerminator([]os.Signal{os.Interrupt})

	term.SetCallback(func(result TerminationResult) {
		if result.FailedOrTimeoutCount != 0 {
			t.Error("FailedOrTimeoutCount should be 0")
			return
		}

		if len(result.Result) != 10 {
			t.Error("Result count should be 10")
			return
		}

		for _, data := range result.Result {
			if data.Error != nil {
				t.Error("Error should be nil")
				return
			}

			if data.Status != SUCCESS {
				t.Error("Status should be SUCCESS")
				return
			}
		}
	})

	result := []string{}
	for i := 0; i < 10; i++ {
		term.Add("app"+strconv.Itoa(i), func(ctx context.Context) error {
			result = append(result, "app"+strconv.Itoa(i))
			return nil
		})
	}

	termInternal := term.(*terminator)
	termInternal.signalChan <- os.Interrupt

	ok := term.Wait(1 * time.Second)
	if !ok {
		t.Error("Wait shouldn't time out")
		return
	}
}
