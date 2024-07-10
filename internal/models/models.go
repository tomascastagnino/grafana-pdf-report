package models

type Panel struct {
	ID int `json:"id"`
	URL string `json:"url"`
	GridPos GridPos `json:"gridPos"`
}

type Dashboard struct {
	Title string `json:"title"`
	Panels []Panel `json:"panels"`
}

type GridPos struct {
    H int `json:"h"`
    W int `json:"w"`
    X int `json:"x"`
    Y int `json:"y"`
}

type PdfRequest struct {
    Panels []Panel `json:"panels"`
}