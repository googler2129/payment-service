package migration

import (
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	upMigration    = "up"
	downMigration  = "down"
	forceMigration = "force"
)

type Migrator interface {
	Up() error
	Down() error
	ForceVersion(version int) error
	MigrateVersion(version uint) error
	Migrate(migrationType string, number string) error
}

type client struct {
	client *migrate.Migrate
}

// InitializeMigrate m, err := migrate.New("file://db/migration", "postgres://localhost:5432/testdb?sslmode=disable&user=Depender&password=password")
func InitializeMigrate(filepath, databaseUrl string) (Migrator, error) {
	m, err := migrate.New(filepath, databaseUrl)
	if err != nil {
		return nil, err
	}

	return &client{
		client: m,
	}, nil
}

func Execute(filepath, databaseUrl, migrationType, number string) {
	migrator, err := InitializeMigrate(filepath, databaseUrl)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
		panic(err)
	}

	migrationErr := migrator.Migrate(migrationType, number)
	if migrationErr != nil {
		log.Fatalf("Failed to execute migration: %v", migrationErr)
		panic(migrationErr)
	}
}

func BuildSQLDBURL(host, port, dbname, username, password string) string {
	return "postgres://" + host + ":" + port + "/" + dbname + "?user=" + username + "&password=" + password + "&sslmode=disable"
}

// Up Apply all Up Migrations
func (c *client) Up() error {
	err := c.client.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Down Apply all Down Migrations
func (c *client) Down() error {
	err := c.client.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// ForceVersion sets a migration version.
func (c *client) ForceVersion(version int) error {
	err := c.client.Force(version)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// MigrateVersion migrated to specific version
func (c *client) MigrateVersion(version uint) error {
	err := c.client.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (c *client) Migrate(migrationType string, number string) error {
	var err error

	switch migrationType {
	case upMigration:
		err = c.Up()
		if err != nil {
			panic(err)
		}
	case downMigration:
		err = c.Down()
		if err != nil {
			panic(err)
		}
	case forceMigration:
		version, parseErr := strconv.Atoi(number)
		if parseErr != nil {
			panic(parseErr)
		}

		err = c.ForceVersion(version)
		if err != nil {
			return err
		}
	default:
		err = c.Up()
		if err != nil {
			panic(err)
		}
	}

	return nil
}
