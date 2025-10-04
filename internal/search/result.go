package search

type Result struct {
	DocID   uint32  `json:"doc_id"`
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Score   float64 `json:"score"`
	Snippet string  `json:"snippet"`
}

type ResultSet []Result

func (r ResultSet) Len() int           { return len(r) }
func (r ResultSet) Less(i, j int) bool { return r[i].Score > r[j].Score }
func (r ResultSet) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
