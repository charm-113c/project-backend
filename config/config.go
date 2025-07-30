/* Package config defines the server's configuration structure
* as well as the names of the configuration files and the function
* to read them.
* IMPORTANT: config files' paths must be defined in environmental variables
* before server startup.
 */
package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config holds the configurations for the server. The CfgPaths are the paths to the .env and .yaml config files,
// and should only be read from environmental variables. The Server vals can be read both from .env and .yaml files.
// CfgPaths is never actually used, it only serves to document the expected key for environmental variables.
type Config struct {
	CfgPaths struct {
		EnvPath string `env:"ENV_PATH" env-default:"config/config.env"`
		CfgPath string `env:"CFG_PATH" env-default:"config/config.yaml"`
	}
	Server struct {
		DevMode  bool   `yaml:"devMode" env:"DEV_MODE" env-default:"true"`
		Port     uint16 `yaml:"port" env:"SRV_PORT" env-default:"7777"`
		Host     string `yaml:"host" env:"SRV_HOST" env-default:"localhost"`
		LogFile  string `yaml:"logFile" env:"LOG_FILE" env-default:"/var/log/junkyard.log"`
		CertFile string `yaml:"certFile" env:"CERT_FILE" env-default:"/some/where/secure"`
		KeyFile  string `yaml:"keyFile" env:"KEY_FILE" env-default:"/some/where/secure"`
	} `yaml:"server"`
	Database struct {
		Host         string `yaml:"host" env:"DB_HOST" env-default:"storage"`
		Type         string `yaml:"type" env:"DB_TYPE" env-default:"postgres"`
		User         string `yaml:"user" env:"DB_USER" env-default:"postgres"`
		Password     string `yaml:"password" env:"DB_PWD" env-default:"password"`
		DBName       string `yaml:"dbName" env:"DB_NAME" env-default:"database"`
		ConnPoolSize int    `yaml:"poolSize" env:"DB_POOL_SIZE" env-default:"20"`
		Port         uint16 `yaml:"port" env:"DB_PORT" env-default:"5432"`
	} `yaml:"database"`
	Router struct {
		// MaxConns     int           `yaml:"maxConns" env:"MAX_CONNS" env-default:"256*1024"` // Let OS decide this
		ReadTimeout  time.Duration `yaml:"readTimeout" env:"READ_TIMEOUT" env-default:"5s"`
		WriteTimeout time.Duration `yaml:"writeTimeout" env:"WRITE_TIMEOUT" env-default:"5s"`
	} `yaml:"router"`
	Logto struct {
		Endpoint  string `yaml:"endpoint" env:"ENDPOINT" env-default:""`
		AppID     string `yaml:"appID" env:"APP_ID" env-default:""`
		AppSecret string `yaml:"appSecret" env:"APP_SECRET" env-default:""`
	} `yaml:"logto"`
}

// LoadConfig reads the config from the config files: values from environmental variables have priority over (i.e. overwrite)
// values from YAML config files, and values from YAML files overwrite the default values (defined in the Config struct)
func LoadConfig(servCfg *Config) error {
	// Load .env file
	envPath := os.Getenv("ENV_PATH")
	// Default to value if envPath is empty
	if envPath == "" {
		envPath = "config.env"
	}
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: could not load .env file %s: %v", envPath, err)
	}

	// Read env vars (if any)
	if err := cleanenv.ReadEnv(servCfg); err != nil {
		err = fmt.Errorf("could not read env variables: %s", err.Error())
		return err
	}

	// Populate CfgPaths should they be needed for debugging or other purposes
	servCfg.CfgPaths.EnvPath = envPath
	cfgPath := os.Getenv("CFG_PATH")
	servCfg.CfgPaths.CfgPath = cfgPath

	// Finally read .yaml config files (if any)
	if cfgPath != "" {
		if err := cleanenv.ReadConfig(cfgPath, servCfg); err != nil {
			err = fmt.Errorf("could not read config file at %s: %w", cfgPath, err)
			return err
		}
	}

	return nil
}

// Validate checks whether the current config values are valid, returning an error if not
func (c *Config) Validate() error {
	if err := ValidateAddress(c.Server.Port, c.Server.Host); err != nil {
		return err
	}
	// if err := ValidateAddress(c.Database.Port, c.Database.Host); err != nil {
	// 	return err
	// }

	// DevMode is loaded into config by cleanenv, so if the input is not a boolean
	// cleanenv will use the default value instead

	// Logfile path should be validated by attempting to create the file
	// But that will be done in startLogger, and would therefore be redundant here

	// TODO: validate DB credentials

	return nil
}

// ValidateAddress takes as input the port and host, and returns an error if they're not valid
func ValidateAddress(port uint16, host string) error {
	if (port < 1024) || (port > 49151) {
		return fmt.Errorf("port number %d is not a valid port", port)
	}

	// Validate host IP
	if (net.ParseIP(host) == nil) && (host != "localhost") {
		return fmt.Errorf("given host address: %s is not valid", host)
	}

	return nil
}
