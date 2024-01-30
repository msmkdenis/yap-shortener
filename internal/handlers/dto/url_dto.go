// Package dto contains data transfer objects.
package dto

type URLResponse struct {
	Result string `json:"result,omitempty"`
}

type URLRequest struct {
	URL string `json:"url,omitempty"`
}

type URLBatchRequest struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}

type URLBatchResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
}

type URLBatchResponseByUserID struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
	DeletedFlag bool   `json:"deleted_flag,omitempty"`
}
