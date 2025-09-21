package indexer

type Redirect struct {
	Title string `xml:"title,attr"`
}

type WikiPage struct {
	Title    string   `xml:"title"`
	ID       int64    `xml:"id"`
	Redirect Redirect `xml:"redirect"`
	Text     string   `xml:"revision>text"`
}

type Parser struct {
	docChan chan<- *Document
	docID   uint32
}
