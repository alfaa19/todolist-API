package main

import (
	"log"
	"net/http"

	activityControllers "github.com/alfaa19/todolist-API/controllers/activityController"
	"github.com/alfaa19/todolist-API/controllers/todoController"
	"github.com/alfaa19/todolist-API/database"
	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()
	r := mux.NewRouter()

	//Activity
	r.HandleFunc("/activity-groups", activityControllers.GetAll).Methods("GET")
	r.HandleFunc("/activity-groups/{id}", activityControllers.GetOne).Methods("GET")
	r.HandleFunc("/activity-groups", activityControllers.Create).Methods("POST")
	r.HandleFunc("/activity-groups/{id}", activityControllers.Update).Methods("PATCH")
	r.HandleFunc("/activity-groups/{id}", activityControllers.Delete).Methods("DELETE")

	//Todo
	r.HandleFunc("/todo-items", todoController.GetAll).Methods("GET")
	r.HandleFunc("/todo-items/{id}", todoController.GetOne).Methods("GET")
	r.HandleFunc("/todo-items", todoController.Create).Methods("POST")
	r.HandleFunc("/todo-items/{id}", todoController.Update).Methods("PATCH")
	r.HandleFunc("/todo-items/{id}", todoController.Delete).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3030", r))

}
