package config

type Config struct {
	Server Server `envPrefix:"SERVER_"`
	AWS    AWS    `envPrefix:"AWS_"`
}

type Server struct {
	Port uint16 `env:"PORT" envDefault:"8080"`

	TLS struct {
		CertPath    string `env:"CERT_PATH,notEmpty"`
		CertKeyPath string `env:"CERT_KEY_PATH,notEmpty"`
	} `envPrefix:"TLS_"`
}

type AWS struct {
	AccessKeyID     string `env:"ACCESS_KEY_ID,notEmpty"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY,notEmpty"`
	Region          string `env:"REGION,notEmpty"`

	DynamoDB struct {
		TableName *string `env:"TABLE_NAME,notEmpty"`
	} `envPrefix:"DYNAMODB_"`
}
