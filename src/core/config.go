package core

// no packages are imported here
// core is a package that contains the application's core logic
// core is imported in other packages

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Env struct {
	PulsarURL                string
	ProjectPort              string
	DBName                   string
	DBPort                   int
	DBPassword               string
	DBUser                   string
	DBHost                   string
	DBSSLMode                string
	MasterDB                 string
	InitOnlyTheseCustomerDBs []string
}

var (
	instance *Env
	once     sync.Once
)

// Load initializes the Env by loading environment variables
func LoadEnv() *Env {

	once.Do(func() {

		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		initDBs := os.Getenv("INIT_ONLY_THESE_CUSTOMER_DBS")
		var initDBList []string
		if initDBs != "" {
			initDBList = strings.Split(initDBs, ",")
		}

		DBPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

		instance = &Env{
			PulsarURL:                os.Getenv("PULSAR_URL"),
			ProjectPort:              os.Getenv("PROJECT_PORT"),
			DBName:                   os.Getenv("DB_NAME"),
			DBPassword:               os.Getenv("DB_PASSWORD"),
			DBUser:                   os.Getenv("DB_USER"),
			DBHost:                   os.Getenv("DB_HOST"),
			DBSSLMode:                os.Getenv("DB_SSLMODE"),
			MasterDB:                 os.Getenv("MASTER_DB_NAME"),
			DBPort:                   DBPort,
			InitOnlyTheseCustomerDBs: initDBList,
		}
	})

	return instance
}
