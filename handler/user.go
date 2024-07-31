package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Leul-Michael/go-auth/global"
	"github.com/Leul-Michael/go-auth/model"
	repository "github.com/Leul-Michael/go-auth/repository/user"
	"github.com/Leul-Michael/go-auth/utils"
)

type User struct {
	Repo *repository.PostgresUserRepo
}

func (h *User) Signup(w http.ResponseWriter, r *http.Request) {
	type userBody struct {
		FirstName   string  `json:"first_name" validate:"required,min=3"`
		LastName    string  `json:"last_name" validate:"required,min=1"`
		Email       string  `json:"email" validate:"required,email"`
		Password    string  `json:"password" validate:"required,min=8"`
		PhoneNumber *string `json:"phone_number" validate:"omitempty,phone"`
	}

	data, err := utils.Decode[userBody](r)
	if err != nil {
		utils.RespondWithError(w, 400, err)
		return
	}

	// validate the request body
	if err := global.Validate.Struct(data); err != nil {
		errs := utils.CustomErrorMessages(err)
		utils.RespondWithError(w, 400, errs)
		return
	}

	// check if user with email exists -> return error
	if count := h.Repo.EmailExists(r.Context(), data.Email); count > 0 {
		utils.RespondWithError(w, 401, "email already exists, please sign in")
		return
	}

	hash, err := utils.HashPassword(data.Password)

	if err != nil {
		utils.RespondWithError(w, 400, "failed to hash password")
		return
	}

	var now = time.Now()
	user := model.User{
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Email:       data.Email,
		Password:    hash,
		PhoneNumber: data.PhoneNumber,
		Role:        "user",
		Base: model.Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := h.Repo.Insert(r.Context(), user); err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("failed to register user: %v", err))
		return
	}

	utils.RespondWithJson(w, 201, struct {
		Message string `json:"message"`
	}{Message: "account created successfully, please sign in"})
}

func (h *User) Login(w http.ResponseWriter, r *http.Request) {
	type userBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	data, err := utils.Decode[userBody](r)
	if err != nil {
		utils.RespondWithError(w, 400, err)
		return
	}

	// validate the request body
	if err := global.Validate.Struct(data); err != nil {
		errs := utils.CustomErrorMessages(err)
		utils.RespondWithError(w, 400, errs)
		return
	}

	sub, err := h.Repo.ComparePassword(r.Context(), data.Email)
	if err != nil {
		utils.RespondWithError(w, 401, "invalid credentials!")
		return
	}

	// check if account is deactivated
	if sub.IsDeactivated {
		utils.RespondWithError(w, 401, "account is deactivated, please contact customer support")
		return
	}

	// check if password matches
	isPasswordCorrect := utils.CheckPasswordHash(string(data.Password), sub.Password)
	if !isPasswordCorrect {
		utils.RespondWithError(w, 401, "invalid credentials!")
		return
	}

	// update last login value
	var now = time.Now()
	if err := h.Repo.UpdateField(r.Context(), sub.Id, "last_login", now); err != nil {
		utils.RespondWithError(w, 401, "invalid credentials!")
		return
	}

	// generate access and refresh tokens
	accessToken, err := utils.GenerateToken(7, sub.Id)
	if err != nil {
		log.Println("failed to generate access token")
		utils.RespondWithError(w, http.StatusInternalServerError, nil)
		return
	}
	refreshToken, err := utils.GenerateToken(24*30, sub.Id) // valid for 30 days
	if err != nil {
		log.Println("failed to generate refresh token")
		utils.RespondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	utils.RespondWithJson(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
