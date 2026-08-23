package dto

import "github.com/golang-jwt/jwt/v5"

type RegistrationInput struct {
	FirstName    string `json:"firstName" schema:"first_name" validate:"required,min=3,max=50"`
	LastName     string `json:"lastName" schema:"last_name" validate:"required,min=3,max=50"`
	Email        string `json:"email" schema:"email" validate:"required,email"`
	HashPassword string `json:"hashPassword" schema:"password" validate:"required,min=6"`
}

type SignIn struct {
	Email    string `json:"email" schema:"email" validate:"required,email"`
	Password string `json:"hashPassword" schema:"password" validate:"required,min=6"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
