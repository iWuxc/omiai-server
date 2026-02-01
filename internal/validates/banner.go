package validates

type BannerListValidate struct {
	Paginate
}

type BannerCreateValidate struct {
	Title     string `json:"title" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required"`
	SortOrder uint   `json:"sort_order"`
	Status    int8   `json:"status" binding:"oneof=0 1"`
	LinkUrl   string `json:"link_url"`
}

type BannerUpdateValidate struct {
	ID        uint64 `json:"id" binding:"required"`
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	SortOrder uint   `json:"sort_order"`
	Status    int8   `json:"status"`
	LinkUrl   string `json:"link_url"`
}

type BannerDeleteValidate struct {
	ID uint64 `uri:"id" binding:"required"`
}
