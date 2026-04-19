package domain

import "time"

type (
	ServiceProperties struct {
		DebugMode      bool
		ServiceName    string         `validate:"required"`
		Timeout        time.Duration  `validate:"required"`
		PoolConnection int            `validate:"required"`
		Redis          RedisConfig    `validate:"required"`
		Database       DatabaseConfig `validate:"required"`
	}

	DatabaseConfig struct {
		IP          string        `validate:"required"`
		Port        int           `validate:"required"`
		User        string        `validate:"required"`
		Password    string        `validate:"required"`
		Name        string        `validate:"required"`
		DialTimeout time.Duration `validate:"required"`
	}

	RedisConfig struct {
		Host     string `validate:"required"`
		Password string
		Database int
	}
)
