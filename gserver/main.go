package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	filePath, _ := filepath.Abs("../db.json")
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &todos)
	if err != nil {
		return nil, err
	}
	return todos.Todos, nil
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := loadTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func main() {
	http.HandleFunc("/todos", todosHandler)

	fmt.Println("Server is running on http://localhost:6001")
	log.Fatal(http.ListenAndServe(":6001", nil))
}
