// Package model contains the model for the application.
package model

// URL represents the URL model.
type URL struct {
	ID            string `db:"id"`
	Original      string `db:"original_url"`
	Shortened     string `db:"short_url"`
	CorrelationID string `db:"correlation_id"`
	UserID        string `db:"user_id"`
	DeletedFlag   bool   `db:"deleted_flag"`
}

// URLStats represents the URL stats.
type URLStats struct {
	Urls  int
	Users int
}
