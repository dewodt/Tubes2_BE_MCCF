package models

// Wikipedia Article Data Structure
type Article struct {
	ID          string `json:"id"`          // Article ID (QXXXXX)
	Title       string `json:"title"`       // Article Title
	Description string `json:"description"` // Article Description
	Thumbnail   string `json:"thumbnail"`   // Article Thumbnail (wikipedia image url)
	URL         string `json:"url"`         // Article URL (wikipedia page url)
}
