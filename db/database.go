package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strings"

	// initialize the mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Config contains the configuration for the database connection.
type Config struct {
	Name string `yaml:"dbname"`
	Host string `yaml:"dbhost"`
	User string `yaml:"dbuser"`
	Pass string `yaml:"dbpass"`
}

// Database containse the sql.DB object.
type Database struct {
	DB *sql.DB
}

// InitDb starts a database connection and sets the DB object in the Database struct.
func (d *Database) InitDb(cfg *Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	d.DB = db
	return nil
}

// ToString returns a string representation of the Config object
func (cfg *Config) ToString() string {
	return fmt.Sprintf("dbname: %s\ndbhost: %s\ndbuser: %s\ndbpass: %s\n", cfg.Name, cfg.Host, cfg.User, cfg.Pass)
}

// SetupDatabase creates the database setup from sql/db.sql.
func (d *Database) SetupDatabase() error {
	f, err := os.Open("./sql/init.sql")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return err
	}
	contents := buf.String()

	queryList := strings.Split(contents, ";")

	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for _, query := range queryList {
		fmt.Print(query)
		if len(query) > 0 {
			_, err := d.DB.Exec(query)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
