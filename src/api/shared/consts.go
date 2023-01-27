package shared

import "time"

const (
	StatusOK       = "success"
	StatusErr      = "failure"
	AccessTokenExp = time.Minute * 5
)
