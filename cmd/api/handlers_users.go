package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/nctqt/casechronicle/internal/database"
	"github.com/nctqt/casechronicle/internal/jsonhelp"
)

type AddUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (cfg *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || !strings.Contains(email, "@") {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "A valid email address is required", nil)
		return
	}

	if len(req.Password) < 8 {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long", nil)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	now := time.Now().UTC()
	newUser, err := cfg.queries.AddUser(r.Context(), database.AddUserParams{
		ID:             uuid.New(),
		Email:          email,
		HashedPassword: string(hashedPassword),
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		// unique constraint error if email already exists
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			jsonhelp.RespondWithError(w, http.StatusConflict, "An account with this email already exists", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}

	// issue JWT token (valid for 24 hours)
	token, err := cfg.makeJWT(newUser.ID, "user", 24*time.Hour)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not issue auth token", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusCreated, AuthResponse{
		User: UserResponse{
			ID:        newUser.ID,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		},
		Token: token,
	})
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	// fetch user record
	user, err := cfg.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// avoid revealing whether the email exists or not
			jsonhelp.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Database error retrieving user", err)
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// issue JWT token (valid for 24 hours)
	role := "user"
	if user.Email == "nctqt@proton.me" {
		role = "admin"
	}

	token, err := cfg.makeJWT(user.ID, role, 24*time.Hour)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not issue auth token", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	})
}

func (cfg *apiConfig) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	user, err := cfg.queries.GetUserByID(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonhelp.RespondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve user", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, user)
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	var req GetUserByEmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Could get user", err)
	}

	user, err := cfg.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonhelp.RespondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve user", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.queries.ListUsers(r.Context())
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not list users", err)
		return
	}

	jsonhelp.RespondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user_id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	err = cfg.queries.DeleteUser(r.Context(), userUUID)
	if err != nil {
		jsonhelp.RespondWithError(w, http.StatusInternalServerError, "Could not delete user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}
