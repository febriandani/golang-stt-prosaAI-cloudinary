package image

type RequestSearchImage struct {
	Keyword string `json:"keyword,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

type ResponseSearchImages struct {
	ID         int64  `json:"id" db:"id"`
	CategoryID int64  `json:"category_id" db:"category_id"`
	Keyword    string `json:"keyword" db:"keyword"`
	URL        string `json:"url" db:"url"`
}
