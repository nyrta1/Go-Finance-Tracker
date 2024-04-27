package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-finance-tracker/internal/models"
	"go-finance-tracker/internal/repository"
	"go-finance-tracker/internal/rest/form"
	"go-finance-tracker/pkg/logger"
	"go-finance-tracker/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

type AuthHandlers struct {
	UserRepo repository.UserRepo
	RoleRepo repository.RoleRepo
}

func NewAuthHandler(userRepo repository.UserRepo, roleRepo repository.RoleRepo) *AuthHandlers {
	return &AuthHandlers{
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
}

func (h *AuthHandlers) Register(ctx *gin.Context) {
	var registerForm form.RegisterInput

	if err := ctx.ShouldBindJSON(&registerForm); err != nil {
		logger.GetLogger().Error("Invalid registration request:", err)
		ctx.JSON(http.StatusBadRequest, &models.CustomResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	if err := validate(registerForm); err != nil {
		logger.GetLogger().Error("Invalid registration request:", err)
		ctx.JSON(http.StatusBadRequest, &models.CustomResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	_, err := h.UserRepo.GetUserByUsername(registerForm.Username)
	if err == nil {
		logger.GetLogger().Error("Account already registered for username:", registerForm.Username)
		ctx.JSON(http.StatusBadRequest, &models.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "The account is already registered",
		})
		return
	}

	var user models.User
	user.Name = registerForm.Name
	user.Surname = registerForm.Surname
	user.Username = registerForm.Username
	user.Email = registerForm.Email
	user.TotalMoney = 0

	role, err := h.RoleRepo.GetByName("USER")
	if err != nil {
		logger.GetLogger().Error("Role not found!")
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  "Role not found!",
		})
		return
	}
	user.Roles = append(user.Roles, *role)

	hashedPassword, err := utils.HashPassword(registerForm.Password)
	if err != nil {
		logger.GetLogger().Error("Unable to hash the password")
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  "Unable to hash the password",
		})
		return
	}
	user.Password = hashedPassword

	if err := h.UserRepo.CreateUser(&user); err != nil {
		logger.GetLogger().Error("Failed to create user:", err)
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	signedToken, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Username)
	if err != nil {
		logger.GetLogger().Error("Failed to generate jwt token:", err)
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(ctx.Writer, &cookie)

	ctx.JSON(http.StatusOK, &models.CustomResponse{
		Status:  http.StatusOK,
		Message: "User registered successfully",
	})
}

func (h *AuthHandlers) Login(ctx *gin.Context) {
	var loginForm form.LoginInput

	if err := ctx.ShouldBindJSON(&loginForm); err != nil {
		logger.GetLogger().Error("Invalid login request:", err)
		ctx.JSON(http.StatusBadRequest, &models.CustomResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	user, err := h.UserRepo.GetUserByUsername(loginForm.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, &models.CustomResponse{
				Status: http.StatusNotFound,
				Error:  err.Error(),
			})
		}
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if !utils.CheckPasswordHash(loginForm.Password, user.Password) {
		logger.GetLogger().Error("Bad credentials for username:", user.Username)
		ctx.JSON(http.StatusBadRequest, &models.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad credentials",
		})
		return
	}

	signedToken, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Username)
	if err != nil {
		logger.GetLogger().Error("Failed to generate jwt token:", err)
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(ctx.Writer, &cookie)

	ctx.JSON(http.StatusOK, &models.CustomResponse{
		Status:  http.StatusOK,
		Message: "User login successful",
	})
}

func (h *AuthHandlers) Logout(ctx *gin.Context) {
	logger.GetLogger().Info("User logout")

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(ctx.Writer, &cookie)

	ctx.JSON(http.StatusOK, &models.CustomResponse{
		Status:  http.StatusOK,
		Message: "User logged out successfully",
	})
}

func (h *AuthHandlers) Profile(ctx *gin.Context) {
	logger.GetLogger().Info("Fetching user profile")

	usernameCtx, exists := ctx.Get("username")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		ctx.JSON(http.StatusUnauthorized, &models.CustomResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	username, ok := usernameCtx.(string)
	if !ok {
		logger.GetLogger().Error("Error while retrieving user ID")
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  "Error while retrieving user ID",
		})
		return
	}

	user, err := h.UserRepo.GetUserByUsername(username)
	if err != nil {
		logger.GetLogger().Error("User does not exist:", err)
		ctx.JSON(http.StatusNotFound, &models.CustomResponse{
			Status:  http.StatusNotFound,
			Message: "User does not exist:",
			Error:   err.Error(),
		})
		return
	}

	logger.GetLogger().Info("User profile fetched successfully")
	ctx.JSON(http.StatusOK, &models.CustomResponse{
		Status:  http.StatusOK,
		Message: "User profile fetched successfully",
		Data:    user,
	})
}
