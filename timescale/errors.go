package timescale

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrSaveEvent      = errors.New("failed to save event to database")
	ErrTransRollback  = errors.New("failed to rollback transaction")
	ErrInvalidEvent   = errors.New("invalid event representation")
)
