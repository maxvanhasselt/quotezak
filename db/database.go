package db

import (
	"database/sql"
	"fmt"
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
	db *sql.DB
}

// InitDb starts a database connection and sets the DB object in the Database struct.
func (d *Database) InitDb(cfg *Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	d.db = db
	return nil
}

// ToString returns a string representation of the Config object
func (cfg *Config) ToString() string {
	return fmt.Sprintf("dbname: %s\ndbhost: %s\ndbuser: %s\ndbpass: %s\n", cfg.Name, cfg.Host, cfg.User, cfg.Pass)
}
