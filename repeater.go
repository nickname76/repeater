package repeater

import (
	"time"
)

// Start repeating fnToCall with passed frequency, frequency can be 0, if so, ticker won't be used.
// Use returned function to stop repeater, this function waits until repeater is stopped.
func StartRepeater(frequency time.Duration, fnToCall func()) (stop func()) {
	if fnToCall == nil {
		panic("fnToCall must not be nil")
	}

	stopCh := make(chan struct{})
	waitCh := make(chan struct{})

	go func() {
		defer close(waitCh)

		if frequency != 0 {
			ticker := time.NewTicker(frequency)

			for {
				select {
				case <-stopCh:
					ticker.Stop()
					return
				case <-ticker.C:
					fnToCall()
				}
			}
		} else {
			for {
				select {
				case <-stopCh:
					return
				default:
					fnToCall()
				}
			}
		}
	}()

	return func() {
		close(stopCh)
		<-waitCh
	}
}
