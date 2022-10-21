package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type AppConfig struct {
	Address     string
	Ports       string
	Port        string
	DB_Driver   string
	DB_Name     string
	DB_Address  string
	DB_Port     string
	DB_Username string
	DB_Password string
	FE_Index    uint
	BE_Index    uint
	QA_Index    uint
	PageLimit   string
	URLFile     string
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {
	var config AppConfig

	config.Address = getEnv("ADDRESS", "http://localhost")
	config.Port = getEnv("PORT", "3000")
	config.Ports = getEnv("PORTS", "443")
	config.DB_Driver = getEnv("DB_DRIVER", "mysql")
	config.DB_Name = getEnv("DB_NAME", "job_vacancies")
	config.DB_Address = getEnv("DB_ADDRESS", "localhost")
	config.DB_Port = getEnv("DB_PORT", "3306")
	config.DB_Username = getEnv("DB_USERNAME", "root")
	config.DB_Password = getEnv("DB_PASSWORD", "")

	config.BE_Index = convertUint(getEnv("BE_INDEX", "2"))
	config.FE_Index = convertUint(getEnv("FE_INDEX", "701"))
	config.QA_Index = convertUint(getEnv("QA_INDEX", "2"))
	config.PageLimit = getEnv("PAGE_LIMIT", "5")
	config.URLFile = getEnv("URL_FILE", "https://google.com")

	fmt.Println(config)
	return &config
}

func PageLimit(data string) string {
	return fmt.Sprintf("/data/%s?limit=%s&page=1", data, GetConfig().PageLimit)
}

func convertUint(data string) uint {
	result, err := strconv.Atoi(data)
	if err != nil {
		fmt.Println("error data index")
	}

	return uint(result)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(value)
		return value
	}

	return fallback
}

func Database() *gorm.DB {
	config := GetConfig()
	switch config.DB_Driver {
	case "mysql":
		return InitMySQL(*config)
	case "postgres":
		return InitPostgreSQL(*config)
	default:
		return nil
	}
}
