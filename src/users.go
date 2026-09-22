package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type userInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var errDuplicateEmail = errors.New("email already exists")

func (app *App) listUsersHandler(w http.ResponseWriter, _ *http.Request) {
	users, err := listUsers(app.DB)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list users")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (app *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input userInput
	if !readJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if message := validateUser(input, true); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return
	}
	user, err := insertUser(app.DB, input)
	if errors.Is(err, errDuplicateEmail) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (app *App) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	var input userInput
	if !readJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if message := validateUser(input, false); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return
	}
	user, err := updateUser(app.DB, id, input)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if errors.Is(err, errDuplicateEmail) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update user")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (app *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	if id == currentUserID(r) {
		writeError(w, http.StatusConflict, "cannot delete the current user")
		return
	}
	result, err := app.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete user")
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateUser(input userInput, passwordRequired bool) string {
	if input.Name == "" {
		return "name is required"
	}
	address, err := mail.ParseAddress(input.Email)
	if err != nil || address.Address != input.Email {
		return "valid email is required"
	}
	if passwordRequired && input.Password == "" {
		return "password is required"
	}
	if input.Password != "" && len(input.Password) < 8 {
		return "password must contain at least 8 characters"
	}
	return ""
}

func ensureInitialUser(db *sql.DB, email, password string) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err := insertUser(db, userInput{Name: "Administrator", Email: strings.ToLower(strings.TrimSpace(email)), Password: password})
	return err
}

func insertUser(db *sql.DB, input userInput) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	createdAt := time.Now().UTC().Format(time.RFC3339)
	result, err := db.Exec(`INSERT INTO users (name, email, password_hash, created_at) VALUES (?, ?, ?, ?)`, input.Name, input.Email, hash, createdAt)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return User{}, errDuplicateEmail
		}
		return User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}
	return User{ID: id, Name: input.Name, Email: input.Email, CreatedAt: createdAt}, nil
}

func updateUser(db *sql.DB, id int64, input userInput) (User, error) {
	var result sql.Result
	var err error
	if input.Password == "" {
		result, err = db.Exec(`UPDATE users SET name = ?, email = ? WHERE id = ?`, input.Name, input.Email, id)
	} else {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return User{}, hashErr
		}
		result, err = db.Exec(`UPDATE users SET name = ?, email = ?, password_hash = ? WHERE id = ?`, input.Name, input.Email, hash, id)
	}
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return User{}, errDuplicateEmail
		}
		return User{}, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return User{}, sql.ErrNoRows
	}
	return userByID(db, id)
}

func userByID(db *sql.DB, id int64) (User, error) {
	var user User
	err := db.QueryRow(`SELECT id, name, email, created_at FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	return user, err
}

func userCredentialsByEmail(db *sql.DB, email string) (User, string, error) {
	var user User
	var hash string
	err := db.QueryRow(`SELECT id, name, email, created_at, password_hash FROM users WHERE email = ?`, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &hash)
	return user, hash, err
}

func listUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT id, name, email, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
