// Package model contains the model for the application.
package model

type URL struct {
	ID            string `db:"id"`
	Original      string `db:"original_url"`
	Shortened     string `db:"short_url"`
	CorrelationID string `db:"correlation_id"`
	UserID        string `db:"user_id"`
	DeletedFlag   bool   `db:"deleted_flag"`
}
