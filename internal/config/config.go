package config

import "os"

type Config struct {
	SQSQueueURL        string
	SQSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	PostgresConnStr    string
	AddressAPIURL      string
}

func LoadConfig() (*Config, error) {
	return &Config{
		SQSQueueURL:        os.Getenv("SQS_QUEUE_URL"),
		SQSRegion:          os.Getenv("SQS_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		PostgresConnStr:    os.Getenv("POSTGRES_CONN_STR"),
		AddressAPIURL:      os.Getenv("ADDRESS_API_URL"),
	}, nil
}
