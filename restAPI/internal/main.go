package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type Todo struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Completed bool `json:"completed"`
}

var todos []Todo

func getTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func addTdos(w http.ResponseWriter, r *http.Request){
	// take value from body and append it in todos
	if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Decode the JSON body into the newTodo variable
    var newTodo Todo
    err := json.NewDecoder(r.Body).Decode(&newTodo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Append the new todo to the todos slice
    todos = append(todos, newTodo)

    // Set the response header and encode the new todo as JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Decode the JSON body into an updatedTodo variable
    var updatedTodo Todo
    err := json.NewDecoder(r.Body).Decode(&updatedTodo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Search for the todo with the matching ID
    for i, todo := range todos {
        if todo.ID == updatedTodo.ID {
            // Update the completed status
            todos[i].Completed = updatedTodo.Completed

            // Set the response header and encode the updated todo as JSON
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(todos[i])
            return
        }
    }

    // If no matching todo is found, return a 404 error
    http.Error(w, "Todo not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r* http.Request){
	if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Decode the JSON body to get the ID of the todo to delete
    var todoToDelete struct {
        ID string `json:"id"`
    }
    err := json.NewDecoder(r.Body).Decode(&todoToDelete)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Search for the todo with the matching ID
    for i, todo := range todos {
        if todo.ID == todoToDelete.ID {
            // Remove the todo from the slice
            todos = append(todos[:i], todos[i+1:]...)

            // Set the response header and status
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "Todo with ID %s has been deleted", todoToDelete.ID)
            return
        }
    }

    // If no matching todo is found, return a 404 error
    http.Error(w, "Todo not found", http.StatusNotFound)

}

func main() {

	PORT := ":8080";

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

	
    // Define routes
    http.HandleFunc("/todos", getTodos)
	http.HandleFunc("/addTodos", addTdos);
	http.HandleFunc("/updateTodo", updateTodo)
	http.HandleFunc("/deleteTodo", deleteTodo)

    log.Println("Listening on localhost:8080")

    log.Fatal(http.ListenAndServe(PORT, nil))
}
