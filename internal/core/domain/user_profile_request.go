package domain

type UserProfileRequest struct {
	UserID    *int    `json:"-" db:"user_id"`
	UserName  *string `json:"username,omitempty" db:"username"`
	UserPhoto *string `json:"photo,omitempty" db:"photo"`
	Gender    *string `json:"gender,omitempty" db:"gender"`
}
