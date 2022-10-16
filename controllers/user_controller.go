package controllers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"time"
	"users-task/configs"
	"users-task/helpers"
	"users-task/models"
	"users-task/repository"
	"users-task/repository/mongodb"
	"users-task/responses"

	"github.com/gin-gonic/gin"
)

var userCollection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()
var userRepo repository.UserRepository

func init() {
	userRepo, _ = mongodb.Connect()
	helpers.CreateFolder("./files")
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: err.Error()}.Log())
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: validationErr.Error()}.Log())
			return
		}

		newUser := &models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			DOB:   user.DOB,
		}

		result, err := userRepo.Create(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: []interface{}{map[string]interface{}{"insertedID": result}}}.Log())
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id, er := strconv.Atoi(c.Param("id"))
		if er != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "param `id` is not a valid user id :" + er.Error()}.Log())
			return
		}

		defer cancel()
		user, err := userRepo.Get(ctx, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: []interface{}{user}}.Log())
	}
}

func EditAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id, er := strconv.Atoi(c.Param("id"))
		if er != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "param `id` is not a valid user id :" + er.Error()}.Log())
			return
		}
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: []interface{}{err.Error()}}.Log())
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: validationErr.Error()}.Log())
			return
		}
		updatedUser, err := userRepo.Update(ctx, id, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: []interface{}{updatedUser}}.Log())
	}
}

func DeleteAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id, er := strconv.Atoi(c.Param("id"))
		if er != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "param `id` is not a valid user id :" + er.Error()}.Log())
			return
		}
		defer cancel()

		err := userRepo.Delete(ctx, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}

		c.JSON(http.StatusOK,
			responses.UserResponse{Status: http.StatusOK, Message: "success", Data: []interface{}{"User successfully deleted!"}},
		)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		users, err := userRepo.List(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}
		c.JSON(http.StatusOK,
			responses.UserResponse{Status: http.StatusOK, Message: "success", Data: []interface{}{users}},
		)
	}
}

func AddFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		id, er := strconv.Atoi(c.Param("id"))
		if er != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "param `id` is not a valid user id :" + er.Error()}.Log())
			return
		}
		userRepo.AddFile(ctx, id)
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}
		log.Println(file.Filename)
		uploadedFileName := file.Filename
		helpers.CreateFolder("./files/" + strconv.Itoa(id))
		file.Filename = helpers.GetFileNameToAdd(strconv.Itoa(id)+"_"+file.Filename, id, 0)
		err = c.SaveUploadedFile(file, file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error()}.Log())
			return
		}

		c.JSON(http.StatusOK,
			responses.UserResponse{Status: http.StatusOK, Message: "success", Data: []interface{}{fmt.Sprintf("file %v uploaded successfully", uploadedFileName)}},
		)
	}
}
