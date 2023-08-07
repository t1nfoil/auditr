package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
	DBName   string `yaml:"dbName"`
}

type siConfig struct {
	SkipErrorCodes []int32 `yaml:"skipErrorCodes"`
}

type auditrConfig struct {
	TenantId       string   `yaml:"tenantId"`
	ClientId       string   `yaml:"clientId"`
	ClientSecret   string   `yaml:"clientSecret"`
	TrimEntryAge   int      `yaml:"trimEntryAge"`
	UpdateInterval int      `yaml:"updateInterval"`
	SiConfig       siConfig `yaml:"signInActivity"`
	Database       dbConfig `yaml:"database"`
}

func (config *auditrConfig) loadConfigYaml() {
	filePath := "/config/auditr.yaml"
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}
}

func (config *auditrConfig) checkConfig() {
	if config.TrimEntryAge < config.UpdateInterval {
		log.Fatalf("config: TrimEntryAge must be greater than UpdateInterval in config/auditr.yaml")
	}
	if config.TenantId == "" {
		log.Fatalf("config: TenantId not set in config/auditr.yaml")
	}
	if config.ClientId == "" {
		log.Fatalf("config: ClientId not set in config/auditr.yaml")
	}
	if config.ClientSecret == "" {
		log.Fatalf("config: ClientSecret not set in config/auditr.yaml")
	}
	if config.TrimEntryAge == 0 {
		log.Fatalf("config: TrimEntryAge not set in config/auditr.yaml")
	}
	if config.UpdateInterval == 0 {
		log.Fatalf("config: UpdateInterval not set in config/auditr.yaml")
	}
	// check postgres config
	if config.Database.Host == "" {
		log.Fatalf("config: postgresConfig.Host not set in config/auditr.yaml")
	}
	if config.Database.Port == 0 {
		log.Fatalf("config: postgresConfig.Port not set in config/auditr.yaml")
	}
	if config.Database.User == "" {
		log.Fatalf("config: postgresConfig.User not set in config/auditr.yaml")
	}
	if config.Database.Password == "" {
		log.Fatalf("config: postgresConfig.Password not set in config/auditr.yaml")
	}
	if config.Database.SSLMode == "" {
		log.Fatalf("config: postgresConfig.SSLMode not set in config/auditr.yaml")
	}
	if config.Database.DBName == "" {
		log.Fatalf("config: postgresConfig.DBName not set in config/auditr.yaml")
	}

}

// create get methods for each of the config variables
func (config *auditrConfig) GetTenantId() string {
	return config.TenantId
}

func (config *auditrConfig) GetClientId() string {
	return config.ClientId
}

func (config *auditrConfig) GetClientSecret() string {
	return config.ClientSecret
}

func (config *auditrConfig) GetTrimEntryAge() int {
	return config.TrimEntryAge
}

func (config *auditrConfig) GetUpdateInterval() int {
	return config.UpdateInterval
}

func (config *auditrConfig) GetSkipErrorCodes() []int32 {
	return config.SiConfig.SkipErrorCodes
}

func (config *auditrConfig) GetDatabaseDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)
	return dsn
}
