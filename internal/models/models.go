package models

type Error struct {
	Error string `json:"error"`
}

type OriginalLink string
type ShortCode string

type ShortCodeRequest struct {
	ShortLink string `json:"short_link"`
}

type ShortCodeResponse struct {
	ShortCode ShortCode `json:"short_code"`
	ShortLink string    `json:"short_link"`
}

type OriginalLinkRequest struct {
	OriginalLink OriginalLink `json:"original_link"`
}

type OriginalLinkResponse struct {
	OriginalLink OriginalLink `json:"original_link"`
}
