package users

import "time"

var (
	RefreshTokenExpiredTime = time.Hour * 24 * 7
	AccessTokenExpiredTime  = time.Minute * 15
)