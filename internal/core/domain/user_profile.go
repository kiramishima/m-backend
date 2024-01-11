package domain

// UserProfile struct
type UserProfile struct {
	UserID    int    `json:"-"`
	UserName  string `json:"username" db:"username"`
	UserPhoto string `json:"photo" db:"photo"`
	Gender    string `json:"gender" db:"gender"`
}
