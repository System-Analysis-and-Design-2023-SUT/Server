package server

import "github.com/pkg/errors"

var ErrNilHealthModule = errors.New("Health module should not be empty")
var ErrNilQueueModule = errors.New("Queue module should not be empty")
