package dto

type URLResponseType struct {
	Result string `json:"result,omitempty"`
}

type URLRequestType struct {
	URL string `json:"url,omitempty"`
}

type URLBatchRequestType struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}

type URLBatchResponseType struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortenedURL  string `json:"short_url,omitempty"`
}
