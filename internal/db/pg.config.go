package db

import "fmt"

type PostgresConfig struct {
	SSLMode string
	Host    string
	Port    string
	Name    string
	User    string
	Pass    string
}

func (c PostgresConfig) String() string {
	return fmt.Sprintf(
		"sslmode=%s host=%s port=%s dbname=%s user=%s password=%s",
		c.SSLMode, c.Host, c.Port, c.Name, c.User, c.Pass,
	)
}
