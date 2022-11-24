package db

import "sync/atomic"

var currentId atomic.Int32
