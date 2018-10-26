package types

import "time"

type Result struct {
	TimeTaken time.Duration
	Success bool
}