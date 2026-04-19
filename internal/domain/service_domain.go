package domain

import "time"

type (
	ServiceProperties struct {
		DebugMode      bool
		ServicePort    int           `validate:"required"`
		ServiceTCPPort int           `validate:"required"`
		ServiceName    string        `validate:"required"`
		Timeout        time.Duration `validate:"required"`
		PoolConnection int           `validate:"required"`
		SmtpConfig     SmtpConfig    `validate:"required"`
	}

	SmtpConfig struct {
		Host     string `validate:"required"`
		Port     int    `validate:"required"`
		Username string `validate:"required"`
		Password string `validate:"required"`
	}
)
