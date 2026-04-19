package config

import (
	"flag"
	"fmt"
	"ordermini-notification-service/internal/domain"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func GetServiceProperties() (svcProperties domain.ServiceProperties) {
	isLocalDev := flag.Bool("local", false, "=(true/false)")
	flag.Parse()

	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	if *isLocalDev {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("error reading transaction file: %w", err))
		}
	} else {
		viper.SetConfigFile(".env.default")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("error reading transaction file: %w", err))
		}
		viper.SetConfigFile(".env.svc")
		if err := viper.MergeInConfig(); err != nil {
			panic(fmt.Errorf("error merging the .env.svc file: %w", err))
		}
	}

	svcProperties = GetEnvServiceProperties()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		svcProperties = GetEnvServiceProperties()
	})
	return svcProperties
}

func GetEnvServiceProperties() (svcProperties domain.ServiceProperties) {
	fmt.Println("Starting Load Env " + time.Now().Format("2006-01-02 15:04:05"))

	svcProperties = domain.ServiceProperties{
		DebugMode:      viper.GetBool("DEBUG_MODE"),
		ServicePort:    viper.GetInt("SERVICE_PORT"),
		ServiceTCPPort: viper.GetInt("SERVICE_TCP_PORT"),
		ServiceName:    viper.GetString("SERVICE_NAME"),
		Timeout:        viper.GetDuration("TIMEOUT"),
		PoolConnection: viper.GetInt("POOL_CONNECTION"),
		SmtpConfig: domain.SmtpConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_AUTH_EMAIL"),
			Password: viper.GetString("SMTP_AUTH_PASSWORD"),
		},
	}

	if err := validator.New().Struct(svcProperties); err != nil {
		panic(err)
	}

	fmt.Println("Finish Load Env " + time.Now().Format("2006-01-02 15:04:05"))
	return svcProperties
}
