package api

import (
	"fmt"
	"os"
	"time"

	"github.com/Jimbo8702/goreservation/db"
	"github.com/Jimbo8702/goreservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email 	 string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User 	*types.User `json:"user"`
	Token 	string 		`json:"token"`
}

type genericResp struct {
	Type	string 	`json:"type"`
	Msg		string
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return ErrInvalidCredentials()
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return ErrInvalidCredentials()
	}
	resp := AuthResponse{
		User: user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id": 		user.ID,
		"email": 	user.Email,
		"expires":  expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}