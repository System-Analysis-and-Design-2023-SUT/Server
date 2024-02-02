package helper

import "github.com/pkg/errors"

var ErrNilMemberlist = errors.New("Helper memberlist should not be nil")
var ErrQueueNotFound = errors.New("Cant get any queue from other nodes")
