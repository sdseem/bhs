package entity

// User -.
type User struct {
	Id           int64  `json:"id"       example:"1"`
	Username     string `json:"username"  example:"username"`
	PasswordHash string `json:"passwordHash"     example:"wuouwou4og34goq"`
}
