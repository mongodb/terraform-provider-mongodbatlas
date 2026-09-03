package acc

import "time"

// WaitUntilGoneOk retries a "still exists" check after delete.
func WaitUntilGoneOk(stillExists func() bool) bool {
	for range 5 {
		if !stillExists() {
			return true
		}
		time.Sleep(3 * time.Second)
	}
	return false
}
