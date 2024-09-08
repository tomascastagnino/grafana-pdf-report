package models

type Panel struct {
	ID      int     `json:"id"`
	URL     string  `json:"url"`
	GridPos GridPos `json:"gridPos"`
	Type    string  `json:"type"`
	Tag     string  `json:"tag"`
}

type Dashboard struct {
	UID       string  `json:"uid"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Panels    []Panel `json:"panels"`
	FolderUid string  `json:"folderUid"`
}

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Options struct {
	Content string `json:"content"`
}

type PdfRequest struct {
	Panels []Panel `json:"panels"`
}
