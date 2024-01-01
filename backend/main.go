package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

var db *sql.DB

func init() {
	// Initialize the database
	var err error
	db, err = sql.Open("sqlite3", "./user.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the user table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			email TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Retrieve users from the database
	users, err := getUsersFromDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching users from the database")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve a specific user by ID
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user ID")
		return
	}

	user, err := getUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create a new user
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request data")
		return
	}

	userID, err := saveUserToDB(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "User could not be created")
		return
	}

	newUser.ID = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Update an existing user by ID
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user ID")
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request data")
		return
	}

	err = updateUserInDB(userID, updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "User could not be updated")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Delete a user by ID
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user ID")
		return
	}

	err = deleteUserFromDB(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "User could not be deleted: %v", err)
		return
	}

	// If deletion is successful, send an empty JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

func getUsersFromDB() ([]User, error) {
	// Retrieve all users from the database
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func getUserByID(userID int) (User, error) {
	// Retrieve a user by ID from the database
	var user User
	err := db.QueryRow("SELECT id, username, email FROM users WHERE id=?", userID).
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func saveUserToDB(user User) (int, error) {
	// Save a new user to the database
	result, err := db.Exec("INSERT INTO users(username, email) VALUES(?, ?)", user.Username, user.Email)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func updateUserInDB(userID int, updatedUser User) error {
	// Update a user in the database
	_, err := db.Exec("UPDATE users SET username=?, email=? WHERE id=?", updatedUser.Username, updatedUser.Email, userID)
	return err
}

func deleteUserFromDB(userID int) error {
	// Delete a user from the database
	_, err := db.Exec("DELETE FROM users WHERE id=?", userID)
	return err
}

func main() {
	r := mux.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	r.HandleFunc("/api/users", GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", GetUser).Methods("GET")
	r.HandleFunc("/api/users", CreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id:[0-9]+}", UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id:[0-9]+}", DeleteUser).Methods("DELETE")

	// Add CORS headers to all routes using CORS middleware
	handler := corsHandler.Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
