package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	db "tarun-kavipurapu/test-go-chat/db/sqlc"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	// "github.com/tarun-kavipurapu/gin-gorm/db/models"
)

// Note:=
//1-->owner
//2->investor

var privateKey []byte = []byte("cEwHkXr2u5x8A/B?D(G+KbPeShVmYq3t6v9y$B&E)H@McQfTjWnZr4u7x!A%C*F-JaN")

func GenerateJWT(user db.User) (string, error) {
	tokenTTL, _ := strconv.Atoi("36000")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"iat":     time.Now().Unix(),
		"eat":     time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(context *gin.Context) error {
	token, err := ExtractJWT(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// /validate Restraunt
func ValidateOwnerRoleJWT(context *gin.Context) error {
	token, err := ExtractJWT(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["user_role"].(float64))
	if ok && token.Valid && userRole == 1 {
		return nil
	}
	return errors.New("invalid owner token provided")
}

// validate Customer Role
func ValidateInvestorRoleJWT(context *gin.Context) error {
	token, err := ExtractJWT(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["user_role"].(float64))
	if ok && token.Valid && userRole == 2 {
		return nil
	}
	return errors.New("invalid auth token provided")
}

//

func CurrentUser(context *gin.Context, store db.Store) (db.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return db.User{}, err
	}
	token, _ := ExtractJWT(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	user, err := store.GetUserById(context, int64(userId))
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func ExtractUserIdFromToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := int64(claims["user_id"].(float64))

	return userId, err
}

func ExtractJWT(context *gin.Context) (*jwt.Token, error) {
	tokenString := ExtractFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func ExtractFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	//will see if i have to extract from the cookie in the future
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
