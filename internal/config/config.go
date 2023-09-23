package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName    string      `json:"serviceName"`
	ServiceAddress string      `json:"servicePort"`
	ServiceID      string      `json:"serviceID"`
	RPCAddress     string      `json:"rpcAddress"`
	Environment    Environment `json:"environment"`

	BuildVer  string
	BuildTime string
	FilePath  string
	RunSince  time.Time

	RedisConfig       RedisConfig `json:"redisConfig"`
	PassthroughConfig PassthroughConfig
}

const logTagConfig = "[Init Config]"

var config *Config

func Init(buildTime, buildVer string) {
	godotenv.Load("conf/.env")

	conf := Config{
		ServiceName:    os.Getenv("SERVICE_NAME"),
		ServiceAddress: os.Getenv("SERVICE_ADDR"),
		ServiceID:      os.Getenv("SERVICE_ID"),
		RPCAddress:     os.Getenv("GPRC_ADDR"),
		RedisConfig: RedisConfig{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		PassthroughConfig: PassthroughConfig{
			SecurityPath: os.Getenv("SP_SECURITY_PATH"),
			AccountPath:  os.Getenv("SP_ACCOUNT_PATH"),
			PaymentPath:  os.Getenv("SP_PAYMENT_PATH"),
		},
		BuildVer:  buildVer,
		BuildTime: buildTime,
		FilePath:  os.Getenv("FILE_PATH"),
	}

	if conf.ServiceName == "" {
		log.Fatalf("%s service name should not be empty", logTagConfig)
	}

	if conf.ServiceAddress == "" {
		log.Fatalf("%s service port should not be empty", logTagConfig)
	}

	envString := os.Getenv("ENVIRONMENT")
	if envString != "dev" && envString != "prod" && envString != "local" {
		log.Fatalf("%s environment must be either local, dev or prod, found: %s", logTagConfig, envString)
	}

	conf.Environment = Environment(envString)

	conf.RunSince = time.Now()
	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
