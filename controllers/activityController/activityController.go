package activityControllers

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
	var activity []models.Activity

	if err := database.DB.Find(&activity).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", "Fail")
		return
	}

	responseSuccess(w, http.StatusOK, activity)

}

func GetOne(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseError(w, http.StatusBadRequest, "Fail", "Bad Request")
		return
	}

	var activity models.Activity

	if err := database.DB.First(&activity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			message := fmt.Sprintf("Activity with ID %d Not Found", id)
			responseError(w, http.StatusNotFound, "Not Found", message)
			return
		} else {
			responseError(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	responseSuccess(w, http.StatusOK, activity)

}

func Create(w http.ResponseWriter, r *http.Request) {
	var activity models.Activity

	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Error", "Invalid request body")
		return
	}

	err = validate.Struct(activity)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			if fieldError.Tag() == "required" {
				message := strings.ToLower(fieldError.Field()) + " cannot be null"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			} else if fieldError.Tag() == "email" {
				message := strings.ToLower(fieldError.Field()) + " format doesnt valid"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			}

		}
	}

	if err := database.DB.Create(&activity).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	responseSuccess(w, http.StatusOK, activity)

}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responseError(w, http.StatusBadRequest, "Fail", "Bad Request")
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			message := fmt.Sprintf("Activity with ID %d Not Found", id)
			responseError(w, http.StatusNotFound, "Not Found", message)
			return
		} else {
			responseError(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	var updateActivity struct {
		ActivityID uint   `gorm:"primaryKey" json:"id"`
		Title      string `gorm:"not null;type:varchar(255)" json:"title" validate:"min=1"`
		Email      string `gorm:"type:varchar(255)" json:"email" validate:"email,min=1"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateActivity)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Error", "Invalid Request Body")

		return
	}
	err = validate.Struct(updateActivity)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			if fieldError.Tag() == "min" {
				message := strings.ToLower(fieldError.Field()) + " cannot be null"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			} else if fieldError.Tag() == "email" {
				message := strings.ToLower(fieldError.Field()) + " format doesnt valid"
				responseError(w, http.StatusBadRequest, "Bad Request", message)
				return
			}

		}

	}

	if err := database.DB.Model(&activity).Updates(updateActivity).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}
	responseSuccess(w, http.StatusOK, activity)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		responseError(w, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	var activity models.Activity

	if err := database.DB.First(&activity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			message := fmt.Sprintf("Activity with ID %d Not Found", id)
			responseError(w, http.StatusNotFound, "Not Found", message)
			return
		} else {
			responseError(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	if err := database.DB.Delete(&activity).Error; err != nil {
		responseError(w, http.StatusInternalServerError, "Error", "Failed to delete activity")
		return
	}

	responseSuccess(w, http.StatusOK, map[string]interface{}{})

}
