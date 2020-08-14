package models

type TokensDTO struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	AccessExpiresIn  int64  `json:"access_expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}
