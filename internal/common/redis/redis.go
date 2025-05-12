package redis

import (
	"fmt"

	rc "github.com/go-redis/redis"
)

type RedisWrapper struct {
	*rc.Client
}

type Option func(r *rc.Options) error

func New(options ...Option) (*RedisWrapper, error) {
	rco := &rc.Options{}

	for _, o := range options {
		err := o(rco)
		if err != nil {
			return nil, err
		}
	}

	cl := rc.NewClient(rco)

	// Validate if configurations are right
	_, err := cl.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &RedisWrapper{cl}, nil
}

func WithAuthentication(password string) Option {
	return func(r *rc.Options) error {
		if password == "" {
			return fmt.Errorf("password can't be empty")
		}
		r.Password = password
		return nil
	}
}

func WithAddr(address string) Option {
	return func(r *rc.Options) error {
		if address == "" {
			return fmt.Errorf("address can't be empty")
		}
		rt, err := rc.ParseURL(address)
		if err != nil {
			return fmt.Errorf("invalid connection string %s", err)
		}
		r.Addr = rt.Addr
		return nil
	}
}
