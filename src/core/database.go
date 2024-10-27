package core

// no packages are imported here
// core is a package that contains the application's core logic
// core is imported in other packages

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() {
	fmt.Println("Initializing database package...")
}

// Initialize connections to active customer databases
func SetRapidorCustomers() error {
	env := LoadEnv()

	masterDSN := getCustomerDSN(env, env.MasterDB)

	masterDB, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to master DB: %w", err)
	}

	// Query active customers
	var customers []RapidorCustomer
	query := masterDB.Table("customer_customer").
		Where("is_active = ?", true).
		Where("database != ?", "rapidor_nodb").
		Select("domain, database")

	if len(env.InitOnlyTheseCustomerDBs) > 0 {
		query = query.Where("database IN ?", env.InitOnlyTheseCustomerDBs)
	}

	if err := query.Scan(&customers).Error; err != nil {
		return fmt.Errorf("failed to fetch active customers: %w", err)
	}

	for _, customer := range customers {
		dsn := getCustomerDSN(env, customer.Database)
		if dsn == "" {
			log.Printf("Invalid DSN for customer %s", customer.Domain)
			continue
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Error connecting to DB for customer %s: %v", customer.Domain, err)
			continue
		}

		// sqlDB, err := db.DB()
		// if err == nil {
		// 	sqlDB.SetMaxIdleConns(10)
		// 	sqlDB.SetMaxOpenConns(100)
		// }

		connMutex.Lock()
		customerConnections[customer.Domain] = db
		connMutex.Unlock()
	}

	return nil
}

func getCustomerDSN(env *Env, database string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		env.DBUser, env.DBPassword, env.DBHost, env.DBPort, database)
}

// // Get the GORM connection for a specific customer
// func getCustomerDB(domain string) (*gorm.DB, error) {
// 	connMutex.Lock()
// 	db, exists := customerConnections[domain]
// 	connMutex.Unlock()
// 	if !exists {
// 		return nil, fmt.Errorf("no active connection for domain %s", domain)
// 	}
// 	return db, nil
// }
