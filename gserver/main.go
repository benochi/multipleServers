package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type Todos struct {
	Todos []Todo `json:"todos"`
}

func loadTodos() ([]Todo, error) {
	var todos Todos

	filePath, err := filepath.Abs("../db.json")
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &todos)
	if err != nil {
		return nil, err
	}
	return todos.Todos, nil
}

func saveTodos(todos []Todo) error {
	todosJSON, err := json.MarshalIndent(Todos{Todos: todos}, "", "    ")
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs("../db.json")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, todosJSON, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodosHandler(w, r)
	case http.MethodPost:
		createTodoHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := loadTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todos, err := loadTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var maxID int
	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}
	newTodo.ID = maxID + 1

	todos = append(todos, newTodo)
	if err := saveTodos(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	http.HandleFunc("/todos", todosHandler)

	fmt.Println("Server is running on http://localhost:6001")
	log.Fatal(http.ListenAndServe(":6001", nil))
}
