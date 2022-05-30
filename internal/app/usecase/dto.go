package usecase

// Correlation is an input DTO for batching
type Correlation struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// OutputBatchItem is an output DTO for batching
type OutputBatchItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// OutputUserLinksListItem is output DTO to represent all User links
type OutputUserLinksListItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
