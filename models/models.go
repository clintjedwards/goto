package models

// Link is a representation of a shortened URL
type Link struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	Created     int64  `json:"created"` // epoch time
	Modified    int64  `json:"modifed"` // last time edited in epoch
	Hits        int64  `json:"hits"`    // number of visits to link
}

// CreateLink is a representation of the createLink payload
type CreateLink struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
