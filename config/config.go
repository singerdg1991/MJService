package config

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Debug                      bool   `mapstructure:"DEBUG"`
	Environment                string `mapstructure:"ENVIRONMENT"`
	Protocol                   string `mapstructure:"PROTOCOL"`
	HostAddress                string `mapstructure:"HOST_ADDRESS"`
	HostUri                    string `mapstructure:"HOST_URI"`
	DockerHostName             string `mapstructure:"DOCKER_HOST_NAME"`
	Port                       int    `mapstructure:"PORT"`
	ApiPrefix                  string `mapstructure:"API_PREFIX"`
	ApiVersion1                string `mapstructure:"API_VERSION_1"`
	ApiVersion2                string `mapstructure:"API_VERSION_2"`
	SigningKey                 string `mapstructure:"SIGNING_KEY"`
	TokenExpiration            int    `mapstructure:"TOKEN_EXPIRATION"`
	RefreshSigningKey          string `mapstructure:"REFRESH_SIGNING_KEY"`
	RefreshTokenExpiration     int    `mapstructure:"REFRESH_TOKEN_EXPIRATION"`
	DatabasePort               int    `mapstructure:"DATABASE_PORT"`
	DatabaseName               string `mapstructure:"DATABASE_NAME"`
	DatabaseHost               string `mapstructure:"DATABASE_HOST"`
	DatabaseUser               string `mapstructure:"DATABASE_USER"`
	DatabasePassword           string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode            string `mapstructure:"DATABASE_SSL_MODE"`
	DatabaseTimeZone           string `mapstructure:"DATABASE_TIMEZONE"`
	Mailable                   bool   `mapstructure:"MAILABLE"`
	GrpcHost                   string `mapstructure:"GRPC_HOST"`
	GrpcPort                   int    `mapstructure:"GRPC_PORT"`
	DatabaseMongoDBHost        string `mapstructure:"DATABASE_MONGODB_HOST"`
	DatabaseMongoDBPort        int    `mapstructure:"DATABASE_MONGODB_PORT"`
	DatabaseMongoDBName        string `mapstructure:"DATABASE_MONGODB_NAME"`
	DatabaseMongoDBUser        string `mapstructure:"DATABASE_MONGODB_USER"`
	DatabaseMongoDBPass        string `mapstructure:"DATABASE_MONGODB_PASS"`
	EventStoreRabbitMQHost     string `mapstructure:"EVENT_STORE_RABBITMQ_HOST"`
	EventStoreRabbitMQPort     int    `mapstructure:"EVENT_STORE_RABBITMQ_PORT"`
	EventStoreRabbitMQUser     string `mapstructure:"EVENT_STORE_RABBITMQ_USER"`
	EventStoreRabbitMQPassword string `mapstructure:"EVENT_STORE_RABBITMQ_PASSWORD"`
	RabbitMQHost               string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPort               int    `mapstructure:"RABBITMQ_PORT"`
	RabbitMQUser               string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPassword           string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMQPanelPort          string `mapstructure:"RABBITMQ_PANEL_PORT"`
	JwtTokenExpiration         int64  `mapstructure:"JWT_TOKEN_EXPIRATION"`
	JwtRefreshTokenExpiration  int64  `mapstructure:"JWT_REFRESH_TOKEN_EXPIRATION"`
	JwtSigningKey              string `mapstructure:"JWT_SIGNING_KEY"`
	OTPCodeLength              int    `mapstructure:"OTP_CODE_LENGTH"`
	OTPCodeExpirationSeconds   int    `mapstructure:"OTP_CODE_EXPIRATION_SECONDS"`
	OTPGRPCHost                string `mapstructure:"OTP_GRPC_HOST"`
	OTPGRPCPort                int    `mapstructure:"OTP_GRPC_PORT"`
	OTPGRPCTimeoutSeconds      int    `mapstructure:"OTP_GRPC_TIMEOUT_SECONDS"`
	OTPTestMode                bool   `mapstructure:"OTP_TEST_MODE"`
	OTPEnable                  bool   `mapstructure:"OTP_ENABLE"`
	MinioEndpoint              string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey             string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey             string `mapstructure:"MINIO_SECRET_KEY"`
	MaxBodySizeLimit           int64  `mapstructure:"MAX_BODY_SIZE_LIMIT"`
	MaxUploadSizeLimit         int64  `mapstructure:"MAX_UPLOAD_SIZE_LIMIT"`
	VAPIDPublicKey             string `mapstructure:"VAPID_PUBLIC_KEY"`
	VAPIDPrivateKey            string `mapstructure:"VAPID_PRIVATE_KEY"`
}

var AppConfig = &Config{}

func Load(path string, configType string, configFileName string) error {
	if path == "" {
		path = "."
	}
	if configType == "" {
		configType = "env"
	}
	if configFileName == "" {
		configFileName = GetEnvPath()
	}
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFileName)

	// read in environment variables that match
	viper.AutomaticEnv()

	// read config from system environment
	elem := reflect.TypeOf(AppConfig).Elem()
	for i := 0; i < elem.NumField(); i++ {
		key := elem.Field(i).Tag.Get("mapstructure")
		value := os.Getenv(key)
		if value == "" {
			break
		}
		viper.Set(key, value)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

	err = viper.Unmarshal(&AppConfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Config unmarshal error: ", err)
		return err
	}

	return nil
}

func LoadDefault() error {
	return Load("", "", "")
}

func GetRootPath() string {
	cmd, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "."
	}
	root := strings.TrimSpace(string(cmd))
	return root
}

func GetEnvPath() string {

	root := GetRootPath()
	abs := fmt.Sprintf("%s/.env", root)
	return abs
}
