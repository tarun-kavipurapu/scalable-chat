package handlers

import (
	"log"
	"net/http"
	db "tarun-kavipurapu/test-go-chat/db/sqlc"
	"tarun-kavipurapu/test-go-chat/types"
	"tarun-kavipurapu/test-go-chat/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store db.Store
}

func NewUserHandler(store db.Store) *UserHandler {
	return &UserHandler{store: store}
}

func (u *UserHandler) CheckUser(ctx *gin.Context, userId int64) {
	_, err := u.store.GetUserById(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User Not Found in the Database"})
	}
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var req types.LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parsing request Error"})
		return
	}

	//get the user by email
	user, err := u.store.GetUserByEmail(ctx, req.UserEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to Fetch the User By Email from DataBase"})
		return
	}
	//compare the password
	err = utils.CheckPassword(req.UserPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is Incorrect"})
		return
	}
	//if no error generate the token and send
	token, err := utils.GenerateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to Generate Token"})
		return
	}

	userDetails := types.UserDetails{
		Id:       int64(user.ID),
		Email:    user.Email,
		Username: user.Username,
	}

	respObject := types.LoginResponse{AccessToken: token, UserDetails: userDetails}
	ctx.JSON(http.StatusOK, types.GenerateResponse(respObject, "Login Successful"))
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	var req types.SignupUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parsing request Error"})
		return
	}
	//hashpassword
	password, err := utils.HashPassword(req.UserPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to hash "})
		return
	}
	log.Println(req.UserEmail)
	log.Println(req.UserName)

	_, err = u.store.CreateUser(ctx, db.CreateUserParams{
		Email:    req.UserEmail,
		Password: password,
		Username: req.UserName,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to Create User in Database"})
		return
	}

	ctx.JSON(http.StatusOK, types.GenerateResponse(nil, "Signup Sucessfull"))

}
