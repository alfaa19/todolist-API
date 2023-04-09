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

	responseSuccess(w, http.StatusCreated, todo)
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

	var updateTodo map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updateTodo)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Error", "Invalid Request Body")

		return
	}
	if len(updateTodo) == 1 {
		for field, value := range updateTodo {
			switch field {
			case "title":
				if err = validate.Var(value, "required,min=1"); err != nil {
					responseError(w, http.StatusBadRequest, "Bad Request", "title cannot be null")
					return
				}
				if err := database.DB.Model(&todo).Update("title", value).Error; err != nil {
					responseError(w, http.StatusInternalServerError, "Error", err.Error())
					return
				}
				responseSuccess(w, http.StatusOK, todo)
				return
			case "activity_group_id":
				if err = validate.Var(value, "required,number,gt=0"); err != nil {
					responseError(w, http.StatusBadRequest, "Bad Request", "activity_group_id cannot be null")
					return
				}
				if err := database.DB.Model(&todo).Update("activity_group_id", value).Error; err != nil {
					responseError(w, http.StatusInternalServerError, "Error", err.Error())
					return
				}
				responseSuccess(w, http.StatusOK, todo)
				return
			case "is_active":
				if err = validate.Var(value, "required,boolean"); err != nil {
					responseError(w, http.StatusBadRequest, "Bad Request", "is_active cannot be null")
					return
				}
				if err := database.DB.Model(&todo).Update("is_active", value).Error; err != nil {
					responseError(w, http.StatusInternalServerError, "Error", err.Error())
					return
				}
				responseSuccess(w, http.StatusOK, todo)
				return
			case "priority":
				if err = validate.Var(value, "required,boolean"); err != nil {
					responseError(w, http.StatusBadRequest, "Bad Request", "priority cannot be null")
					return
				}
				if err := database.DB.Model(&todo).Update("priority", value).Error; err != nil {
					responseError(w, http.StatusInternalServerError, "Error", err.Error())
					return
				}
				responseSuccess(w, http.StatusOK, todo)
				return
			default:
				responseError(w, http.StatusBadRequest, "Bad Request", "Invalid field name")
				return
			}
		}
	}

	var updateFields []string
	for field, value := range updateTodo {
		switch field {
		case "title":
			if err = validate.Var(value, "required,min=1"); err != nil {
				responseError(w, http.StatusBadRequest, "Bad Request", "title cannot be null")
				return
			}
			updateFields = append(updateFields, "title")
		case "activity_group_id":
			if err = validate.Var(value, "required,number,gt=0"); err != nil {
				responseError(w, http.StatusBadRequest, "Bad Request", "Invalid activity_group_id format")
				return
			}
			updateFields = append(updateFields, "activity_group_id")
		case "is_active":
			if err = validate.Var(value, "required,boolean"); err != nil {
				responseError(w, http.StatusBadRequest, "Bad Request", "Invalid is_active format")
				return
			}
			updateFields = append(updateFields, "activity_group_id")
		case "priority":
			if err = validate.Var(value, "required,min=1"); err != nil {
				responseError(w, http.StatusBadRequest, "Bad Request", "Invalid priority format")
				return
			}
			updateFields = append(updateFields, "activity_group_id")
		default:
			responseError(w, http.StatusBadRequest, "Bad Request", "Invalid field name")
			return
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
