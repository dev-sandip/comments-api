package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/dev-sandip/comments-api/internals/utils"
)

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// APIResponse represents a standard API response.
type APIResponse struct {
	Message string      `json:"message" example:"Login successful"`
	Data    interface{} `json:"data,omitempty"`
}

// UserResponse represents public user information.
type UserResponse struct {
	ID    string `json:"id" example:"123"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// SignUpHandler godoc
//
// @Summary      Register a new user
// @Description  Create a new user account and set JWT authentication cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body CreateUserParams true "User registration details"
// @Success      201 {object} APIResponse "User created successfully"
// @Header       201 {string} Set-Cookie "Authentication cookie"
// @Failure      400 {object} APIResponse "Bad request"
// @Failure      409 {object} APIResponse "User already exists"
// @Failure      500 {object} APIResponse "Internal server error"
// @Router       /signup [post]
func (h *handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var params CreateUserParams

	// Read request body.
	if err := utils.Read(r, &params); err != nil {
		log.Printf("Error reading request body: %v", err)
		utils.Fail(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Check whether user already exists.
	data, err := h.service.GetUserByEmail(r.Context(), params.Email)

	if err != nil && !errors.Is(err, ErrUserNotFound) {
		log.Printf("Error checking user: %v", err)
		utils.Fail(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if err == nil && data.Email != "" {
		utils.Fail(w, http.StatusConflict, "User already exists")
		return
	}

	// Hash password.
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.Fail(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create user.
	newUser, err := h.service.CreateUser(r.Context(), CreateUserParams{
		Name:     params.Name,
		Email:    params.Email,
		Password: hashedPassword,
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.Fail(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate JWT.
	token, err := GenerateToken(newUser.ID.String())
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.Fail(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set authentication cookie.
	SetAuthCookie(w, token)

	// Return public user information.
	utils.Success(
		w,
		http.StatusCreated,
		"User created successfully",
		UserResponse{
			ID:    newUser.ID.String(),
			Name:  newUser.Name,
			Email: newUser.Email,
		},
	)
}

// LoginHandler godoc
//
// @Summary      Login user
// @Description  Authenticate user and set JWT authentication cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginParams true "Login credentials"
// @Success      200 {object} APIResponse "Login successful"
// @Header       200 {string} Set-Cookie "Authentication cookie"
// @Failure      400 {object} APIResponse "Bad request"
// @Failure      401 {object} APIResponse "Invalid email or password"
// @Failure      500 {object} APIResponse "Internal server error"
// @Router       /login [post]
func (h *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var params LoginParams

	// Read request body.
	if err := utils.Read(r, &params); err != nil {
		log.Printf("Error reading request body: %v", err)
		utils.Fail(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Find user by email.
	user, err := h.service.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			utils.Fail(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		log.Printf("Error finding user: %v", err)
		utils.Fail(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Verify password.
	if !utils.CheckPasswordHash(params.Password, user.Password) {
		utils.Fail(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT.
	token, err := GenerateToken(user.ID.String())
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.Fail(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set authentication cookie.
	SetAuthCookie(w, token)

	// Return public user information.
	utils.Success(
		w,
		http.StatusOK,
		"Login successful",
		UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	)
}
