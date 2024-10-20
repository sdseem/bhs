package entity

// Asset -.
type Asset struct {
	Id          int64  `json:"id"       example:"1"`
	Name        string `json:"name"  example:"asset 1"`
	Description string `json:"description"     example:"asset description"`
	Price       string `json:"price"  example:"100.00"`
}
