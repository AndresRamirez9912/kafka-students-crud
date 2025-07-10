package api

import "time"

const (
	// how long server waits the pong before close the connection
	PONG_WAIT = 10 * time.Second

	// PINT_INTERVAL  must be lower than pong wait (in this case is 90%)
	PINT_INTERVAL = (PONG_WAIT * 9) / 10
)
