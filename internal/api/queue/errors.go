package queue

import "github.com/pkg/errors"

var ErrNilQueueRepo = errors.New("Queue repository should not be nil")
var ErrNilQueueService = errors.New("Queue service should not be nil")
