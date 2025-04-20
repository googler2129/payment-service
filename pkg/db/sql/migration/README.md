# Database Migration Package

## Overview

The migration package provides a simplified interface for managing PostgreSQL database migrations using the [golang-migrate](https://github.com/golang-migrate/migrate/v4) library. It allows for applying migrations up or down, forcing a specific version, and migrating to a target version.

## Features

- **Up Migration**: Apply all pending migrations
- **Down Migration**: Rollback all migrations
- **Force Version**: Set a migration version directly without running migrations
- **Target Version Migration**: Migrate to a specific version
- **Simple API**: Wrapper around golang-migrate with simplified error handling

## Usage Examples

### Initialize Migration Client

```go
// Initialize a migration client
migrator, err := migration.InitializeMigrate(
    "file:///path/to/migrations",
    "postgres://hostname:5432/dbname?user=username&password=password&sslmode=disable",
)
if err != nil {
    log.Fatalf("Failed to initialize migrator: %v", err)
}
```

### Run Migrations

```go
// Apply all pending up migrations
err := migrator.Up()
if err != nil {
    log.Fatalf("Migration failed: %v", err)
}
```

### Rollback Migrations

```go
// Rollback all migrations
err := migrator.Down()
if err != nil {
    log.Fatalf("Rollback failed: %v", err)
}
```

### Force Specific Version

```go
// Force a specific migration version (useful for troubleshooting)
err := migrator.ForceVersion(3)
if err != nil {
    log.Fatalf("Failed to force version: %v", err)
}
```

### Migrate to Specific Version

```go
// Migrate to a specific version
err := migrator.MigrateVersion(5)
if err != nil {
    log.Fatalf("Failed to migrate to version: %v", err)
}
```

### Execute Migrations From Command Line

```go
// Execute migrations with type and number
migration.Execute(
    "file:///path/to/migrations",
    "postgres://hostname:5432/dbname?user=username&password=password&sslmode=disable",
    "up",   // Migration type: "up", "down", or "force"
    "",     // Number (used with "force" type)
)
```

### Build PostgreSQL Database URL

```go
// Build a PostgreSQL connection URL
dbURL := migration.BuildSQLDBURL(
    "localhost",    // Host
    "5432",         // Port
    "mydatabase",   // Database name
    "myuser",       // Username
    "mypassword",   // Password
)
```

## Migration Types

The package supports the following migration types:

- `up`: Apply all pending migrations
- `down`: Rollback all migrations
- `force`: Force a specific migration version

## Practical Example

Here's a complete example of using the migration package in a project:

```go
package main

import (
    "log"
    
    "github.com/mercor/payment-service/pkg/db/sql/migration"
)

func main() {
    // Build the database URL
    dbURL := migration.BuildSQLDBURL(
        "localhost",
        "5432",
        "application_db",
        "app_user",
        "app_password",
    )
    
    // Path to migration files
    migrationPath := "file://migrations"
    
    // Initialize the migrator
    migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
    if err != nil {
        log.Fatalf("Failed to initialize migrator: %v", err)
    }
    
    // Apply all pending migrations
    err = migrator.Up()
    if err != nil {
        log.Fatalf("Failed to apply migrations: %v", err)
    }
    
    log.Println("Migrations applied successfully!")
}
```

## Important Notes

1. Migration files should follow the format `version_description.up.sql` and `version_description.down.sql`
2. The package handles `migrate.ErrNoChange` errors and returns nil when no migrations are needed
3. Ensure your database connection string is correctly formatted for PostgreSQL
4. Make sure all migration files are accessible from the path provided

## Error Handling

The package handles common migration errors including:

- Connection errors during initialization
- Migration errors during execution
- No changes required (returns nil instead of error)

## Dependencies

This package relies on:
- github.com/golang-migrate/migrate/v4
- github.com/golang-migrate/migrate/v4/database/postgres
- github.com/golang-migrate/migrate/v4/source/file 