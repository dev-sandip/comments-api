package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/dev-sandip/comments-api/internals/utils"
)

var ErrUserNotFound = errors.New("user not found")

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var params CreateUserParams

	if err := utils.Read(r, &params); err != nil {
		log.Printf("Error reading request body: %v", err)
		utils.Fail(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Check whether the user already exists.
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

	utils.Success(
		w,
		http.StatusCreated,
		"User created successfully",
		newUser,
	)
}
