// Package dto contains data transfer objects.
package dto

// URLResponse represents URL response.
type URLResponse struct {
	Result string `json:"result,omitempty"`
}

// URLRequest represents URL request.
type URLRequest struct {
	URL string `json:"url,omitempty"`
}

// URLBatchRequest represents URL batch request.
type URLBatchRequest struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}

// URLBatchResponse represents URL batch response.
type URLBatchResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
}

// URLBatchResponseByUserID represents URL batch response by user ID.
type URLBatchResponseByUserID struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
	DeletedFlag bool   `json:"deleted_flag,omitempty"`
}
