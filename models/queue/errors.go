package models

import "github.com/pkg/errors"

var ErrKeyExist = errors.New("Duplicate key in queue")
var ErrEmptyList = errors.New("Queue is empty")
var ErrParseData = errors.New("Can not parse input data")
var ErrKeyNotFound = errors.New("Key not found")
var ErrObjectNotFound = errors.New("Object not found")

var ErrSubscriberExist = errors.New("You already subscribed")
