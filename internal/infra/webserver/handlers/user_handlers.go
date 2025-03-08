package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JMKobayashi/Basic-API-Server/internal/dto"
	"github.com/JMKobayashi/Basic-API-Server/internal/entity"
	"github.com/JMKobayashi/Basic-API-Server/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db}
}

// Generate JWT godoc
// @Summary Generate JWT Token
// @Description Generate JWT Token
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body dto.GetJWTInput true "user credentials"
// @Success 200 {object} dto.GetJWTOutput
// @Failure 404
// @Failure 500 {object} Error
// @Router /user/token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, token, err := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken := dto.GetJWTOutput{AccessToken: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body dto.CreateUserInput true "user request"
// @Success 201
// @Failure 500 {object} Error
// @Router /user [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
