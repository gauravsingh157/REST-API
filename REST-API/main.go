package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// User struct
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users  = []User{} // In-memory storage
	nextID = 1        // ID counter
	mutex  sync.Mutex // For thread-safe operations
)

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Handler to manage all users (GET, POST)
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUsers(w)
	case "POST":
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler to manage a single user (GET, PUT, DELETE)
func userHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getUser(w, id)
	case "PUT":
		updateUser(w, r, id)
	case "DELETE":
		deleteUser(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /users
func getUsers(w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST /users
func createUser(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.ID = nextID
	nextID++
	users = append(users, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GET /users/{id}
func getUser(w http.ResponseWriter, id int) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// PUT /users/{id}
func updateUser(w http.ResponseWriter, r *http.Request, id int) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {
		if user.ID == id {
			var updatedUser User
			if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			updatedUser.ID = id
			users[i] = updatedUser

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// DELETE /users/{id}
func deleteUser(w http.ResponseWriter, id int) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}
