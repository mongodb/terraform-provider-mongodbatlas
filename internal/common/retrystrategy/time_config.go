package retrystrategy

import "time"

type TimeConfig struct {
	Timeout    time.Duration
	MinTimeout time.Duration
	Delay      time.Duration
}
