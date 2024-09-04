package entity

type Response struct {
	Items []Item `json:"items"`
}

type Item struct {
	Snippet SnippetData `json:"snippet"`
}

type SnippetData struct {
	Title      string        `json:"title"`
	Thumbnails ThumbnailType `json:"thumbnails"`
}

type ThumbnailType struct {
	Default  ThumbnailBody `json:"default"`
	Medium   ThumbnailBody `json:"medium"`
	High     ThumbnailBody `json:"high"`
	Standard ThumbnailBody `json:"standard"`
	Maxres   ThumbnailBody `json:"maxres"`
}

type ThumbnailBody struct {
	Url    string `json:"url"`
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
}
