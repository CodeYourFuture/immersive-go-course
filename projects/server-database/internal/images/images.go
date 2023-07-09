package images

import "errors"

var (
	ImagesNotFound            = errors.New("images not found")
	ImagesUniqueCodeViolation = errors.New("unique code violation")
)

type Image struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	AltText    string `json:"alt_text"`
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
	CreatedAt  string `json:"createdAt"`
}

type CreateImagePayload struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

func (c *CreateImagePayload) Validate() error {
	return nil
}

type CreateImageResponse struct {
	ID string `json:"id"`
}
