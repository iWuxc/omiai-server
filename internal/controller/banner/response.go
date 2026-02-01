package banner

type BannerResponse struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkUrl  string `json:"link_url"`
}
