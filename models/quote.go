package models

import (
	"database/sql"
	"fmt"
)

// Quote represents a quote, the person that said it and the date it was uttered.
type Quote struct {
	Name     string
	Quote    string
	Owner    string
	Date     string
	Category string
}

// NewQuote constructs a quote
func NewQuote(name string, quote string, owner string, date string, category string) *Quote {
	q := new(Quote)
	q.Name = name
	q.Quote = quote
	q.Owner = owner
	q.Date = date
	q.Category = category
	return q
}

// Save saves a quote to the database
func (q *Quote) Save(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = db.Query("INSERT INTO quote (quote_name, quote, owner, date, category) VALUES (?, ?, ?, ?, ?);", q.Name, q.Quote, q.Owner, q.Date, q.Category)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Update updates the quote in the database
func (q *Quote) Update(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = db.Query("UPDATE quote SET quote_name = ?, quote = ?, owner = ?, date = ? WHERE quote_name = ?", q.Name, q.Quote, q.Owner, q.Date, q.Name)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the quote from the database.
func (q *Quote) Delete(db *sql.DB, name string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = db.Query("DELETE FROM quote WHERE quote_name LIKE ?", name)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetQuoteByName retrieves a quote from the database and returns a quote object
func GetQuoteByName(db *sql.DB, name string) (*Quote, error) {

	qu := &Quote{}

	row := db.QueryRow("SELECT quote_name, quote, owner, date FROM quote WHERE quote_name LIKE ?;", name)

	err := row.Scan(&qu.Name, &qu.Quote, &qu.Owner, &qu.Date)
	if err != nil {
		return nil, err
	}
	return qu, nil
}

// GetRandomQuote selects a random quote from the database and returns it.
func GetRandomQuote(db *sql.DB) (*Quote, error) {

	qu := &Quote{}

	rows := db.QueryRow("SELECT quote_name, quote, owner, date FROM quote ORDER BY RAND() LIMIT 1")

	err := rows.Scan(&qu.Name, &qu.Quote, &qu.Owner, &qu.Date)

	if err != nil {
		return nil, err
	}
	return qu, nil
}

// GetRandomOwner selects a random quote by the owner from the database and returns it.
func GetRandomOwner(db *sql.DB, owner string) (*Quote, error) {

	qu := &Quote{}

	rows := db.QueryRow("SELECT quote_name, quote, owner, date FROM quote WHERE owner LIKE ? ORDER BY RAND() LIMIT 1", owner)

	err := rows.Scan(&qu.Name, &qu.Quote, &qu.Owner, &qu.Date)

	if err != nil {
		return nil, err
	}
	return qu, nil
}

// GetRandomCategory selects a random quote from the category from the database and returns it.
func GetRandomCategory(db *sql.DB, cat string) (*Quote, error) {

	qu := &Quote{}

	rows := db.QueryRow("SELECT quote_name, quote, owner, date FROM quote WHERE category LIKE ? ORDER BY RAND() LIMIT 1", cat)

	err := rows.Scan(&qu.Name, &qu.Quote, &qu.Owner, &qu.Date)

	if err != nil {
		return nil, err
	}
	return qu, nil
}

func (q *Quote) String() string {
	return fmt.Sprintf("\"%s\" - %s %s", q.Quote, q.Owner, q.Date)
}

func (q *Quote) Data() string {
	return fmt.Sprintf("Name: %s, Quote: %s, By: %s, Date: %s", q.Name, q.Quote, q.Owner, q.Date)
}
