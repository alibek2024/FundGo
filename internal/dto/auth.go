package dto

import "github.com/golang-jwt/jwt/v5"

type RegistrationInput struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	HashPassword string `json:"hashPassword"`
}

type SignIn struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"hashPassword"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
