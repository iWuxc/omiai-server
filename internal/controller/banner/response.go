package banner

type BannerResponse struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkUrl  string `json:"link_url"`
}
