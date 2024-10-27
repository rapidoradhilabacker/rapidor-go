package core

import (
	"sync"

	"gorm.io/gorm"
)

// Define the customer struct to store domain and database info
type RapidorCustomer struct {
	Domain   string
	Database string
}

// Store active database connections for each customer
var customerConnections = make(map[string]*gorm.DB)
var connMutex sync.Mutex
