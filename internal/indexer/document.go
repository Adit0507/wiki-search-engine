package indexer

type Document struct {
	ID      uint32         `json:"id"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	URL     string         `json:"url"`
	Terms   map[string]int `json:"terms"`
	Length  int            `json:"length"`
}
