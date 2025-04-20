# Package db

## Overview
The db package provides an abstraction layer for database interactions, including connection management, query execution, and transaction handling, simplifying access to various databases.

## Key Components
- Connection Management: Open and close connections efficiently.
- Query Execution: Helper functions to run queries and parse results.
- Transaction Support: Begin, commit, and rollback transactions reliably.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/mercor/payment-service/pkg/db"
)

func main() {
	database, err := db.Connect(db.Config{
		// DSN and driver configurations
	})
	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}
	defer database.Close()
	fmt.Println("Connected to database successfully")
}
~~~

## Notes
- Compatible with multiple database backends.
