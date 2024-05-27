package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	sync.RWMutex
	users map[string]string
	todos map[string][]Todo
}

var store = Store{
	users: make(map[string]string),
	todos: make(map[string][]Todo),
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	store.Lock()
	store.users[creds.Username] = string(hashedPassword)
	store.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	store.RLock()
	hashedPassword, exists := store.users[creds.Username]
	store.RUnlock()

	if !exists || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	todo.ID = len(store.todos[username]) + 1
	todo.CreatedAt = time.Now()

	store.Lock()
	store.todos[username] = append(store.todos[username], todo)
	store.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	store.RLock()
	todos := store.todos[username]
	store.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	store.Lock()
	defer store.Unlock()
	for i, todo := range store.todos[username] {
		if fmt.Sprintf("%d", todo.ID) == id {
			store.todos[username][i] = updatedTodo
			json.NewEncoder(w).Encode(updatedTodo)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	store.Lock()
	defer store.Unlock()
	for i, todo := range store.todos[username] {
		if fmt.Sprintf("%d", todo.ID) == id {
			store.todos[username] = append(store.todos[username][:i], store.todos[username][i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func getUsernameFromRequest(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return ""
	}

	tokenStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return ""
	}
	if !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}
	return claims.Username
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/signin", SignIn).Methods("POST")
	r.HandleFunc("/todos", CreateTodo).Methods("POST")
	r.HandleFunc("/todos", GetTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
