package redis

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

var (
	ErrNil                    = redis.Nil
	ErrInvalidExpireParameter = errors.New("expire must be larger than 0")
)
