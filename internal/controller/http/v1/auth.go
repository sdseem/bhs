package v1

import (
	"bhs/internal/entity"
	"bhs/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"

	"bhs/pkg/logger"
)

type authRoutes struct {
	a usecase.Auth
	l logger.Interface
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.Auth, l logger.Interface) {
	r := &authRoutes{a, l}

	h := handler.Group("/auth")
	{
		h.POST("/registration", r.registration)
		h.POST("/auth", r.auth)
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

// @Summary     Register
// @Description Register a user
// @ID          register
// @Tags  	    auth
// @Accept      json
// @Produce     json
// @Param       request body authRequest
// @Success     200 {object} authResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /auth/registration [post]
func (r *authRoutes) registration(c *gin.Context) {
	var request authRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - registration")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	hash, err := r.a.HashPassword(request.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - registration")
		errorResponse(c, http.StatusBadRequest, "cannot hash password")

		return
	}

	user := entity.User{
		Username:     request.Username,
		PasswordHash: hash,
	}
	token, err := r.a.Register(c, user)
	if err != nil {
		r.l.Error(err, "http - v1 - registration")
		errorResponse(c, http.StatusBadRequest, "registration failed")

		return
	}

	c.JSON(http.StatusOK, authResponse{token})
}

// @Summary     Authentication
// @Description Authenticate a user
// @ID          authenticate
// @Tags  	    auth
// @Accept      json
// @Produce     json
// @Param       request body authRequest
// @Success     200 {object} authResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /auth/auth [post]
func (r *authRoutes) auth(c *gin.Context) {
	var request authRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - auth")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	token, err := r.a.Authenticate(c, request.Username, request.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - registration")
		errorResponse(c, http.StatusBadRequest, "authentication failed")

		return
	}

	c.JSON(http.StatusOK, authResponse{token})
}
