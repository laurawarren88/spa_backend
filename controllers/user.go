package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	userCollection *mongo.Collection
}

func NewUserController(collection *mongo.Collection) *UserController {
	return &UserController{userCollection: collection}
}

func (uc *UserController) AddUser(ctx *gin.Context) {
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
}

func (uc *UserController) SignupUser(ctx *gin.Context) {
}

func (uc *UserController) ForgotPassword(ctx *gin.Context) {
}

func (uc *UserController) LogoutUser(ctx *gin.Context) {
}
