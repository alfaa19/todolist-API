package todoController

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alfaa19/todolist-API/database"
	"github.com/alfaa19/todolist-API/helpers"
	"github.com/alfaa19/todolist-API/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var responseSuccess = helpers.ResponseSuccess
var responseError = helpers.ResponseError
var validate = validator.New()

func GetAll(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo []models.Todo
	if val, ok := params["activity_group_id"]; ok {
		activity_id, err := strconv.Atoi(val)
		if err != nil {
			responseError(w, http.StatusBadRequest, "Bad Request", "Bad Request")
			return
		}
		if err := database.DB.Where("activity_group_id <> ?", activity_id).Find(&todo).Error; err != nil {
			responseError(w, http.StatusInternalServerError, "Error", "Fail")
			return
		}

		responseSuccess(w, http.StatusOK, todo)
		return
	}

	if err := database.DB.Find(&todo).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", "Fail")
		return
	}

	responseSuccess(w, http.StatusOK, todo)

}

func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseError(w, http.StatusBadRequest, "Fail", "Bad Request")
		return
	}

	var todo models.Todo

	if err := database.DB.First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			message := fmt.Sprintf("Todo with ID %d Not Found", id)
			responseError(w, http.StatusNotFound, "Not Found", message)
			return
		} else {
			responseError(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	responseSuccess(w, http.StatusOK, todo)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Error", "Invalid request body")
		return
	}

	err = validate.Struct(todo)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			if fieldError.Tag() == "required" {
				message := strings.ToLower(fieldError.Field()) + " cannot be null"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			} else {
				message := strings.ToLower(fieldError.Field()) + " type missmatch"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			}

		}
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	responseSuccess(w, http.StatusOK, todo)
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responseError(w, http.StatusBadRequest, "Fail", "Bad Request")
		return
	}

	var todo models.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			message := fmt.Sprintf("Todo with ID %d Not Found", id)
			responseError(w, http.StatusNotFound, "Not Found", message)
			return
		} else {
			responseError(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	var updateTodo struct {
		TodoID          uint   `gorm:"primaryKey" json:"id"`
		ActivityGroupID uint   `json:"activity_group_id" validate:"number,gt=0"`
		Title           string `gorm:"type:varchar(255)" json:"title" validate:"min=1"`
		IsActive        bool   `gorm:"type:bool" json:"is_active" validate:"boolean"`
		Priority        string `gorm:"type:varchar(55)" json:"priority" validate:"min=1"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateTodo)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Error", "Invalid Request Body")

		return
	}
	err = validate.Struct(updateTodo)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			if fieldError.Tag() == "min" {
				message := strings.ToLower(fieldError.Field()) + " cannot be null"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			} else {
				responseError(w, http.StatusBadRequest, "Bad Request", err.Error())
				return
			}

		}

	}

	if err := database.DB.Model(&todo).Updates(updateTodo).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}
	responseSuccess(w, http.StatusOK, todo)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseError(w, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	var todo models.Todo

	if err := database.DB.First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			message := fmt.Sprintf("Todo with ID %d Not Found", id)
			responseError(w, http.StatusNotFound, "Not Found", message)
			return
		} else {
			responseError(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	if err := database.DB.Delete(&todo).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", "Failed to delete activity")
		return
	}

	responseSuccess(w, http.StatusOK, map[string]interface{}{})
}
