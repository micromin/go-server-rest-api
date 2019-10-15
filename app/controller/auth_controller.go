package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gos/app/auth"
	"gos/app/models"
	"net/http"
	"time"
)

// swagger:operation POST /api/aut/login Login
//
// Login holds the functionality for login
// ---
// produces:
// - application/json
// parameters:
// - name: body
//   in: body
//   description: the login obj
//   schema:
//    $ref: '#/definitions/LoginRequest'
// responses:
//  '200':
//    description: successful operation
//    schema:
//     $ref: '#/definitions/Response'
//  '400':
//    description: invalid request
//    schema:
//     $ref: '#/definitions/Response'
//  '500':
//    description: internal server error
//    schema:
//     $ref: '#/definitions/LoginResponse'
func (c *AppController) Login(ctx *gin.Context) {
	request := new(models.LoginRequest)
	if err := ctx.BindJSON(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", err))
		return
	}

	user, err := c.appRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, getErrorResponse("email or password is incorrect", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		user.FailedLoginAttempt++
		_, err := c.appRepo.UpdateUser(ctx, *user)
		if err != nil {
			fmt.Println(errors.Wrap(err, "failed to update FailedLoginAttempt"))
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, getErrorResponse("email or password is incorrect", errors.New("email or password is incorrect")))
		return
	}

	expiresAt := time.Now().UTC().Add(auth.AccessTokenExpirationMinutes * time.Minute).Unix()
	claims := &auth.Claims{
		UserId: user.UserId,
		Email:  user.Email,
		Name:   user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(c.auth.GetJWTKey())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, getErrorResponse("email or password is incorrect", err))
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Message: "successfully logged in user",
		Data: models.LoginResponse{
			Token:     tokenString,
			ExpiresAt: expiresAt,
		},
	})
}

// swagger:operation POST /api/aut/register Register
//
// Register holds the functionality for registration
// ---
// produces:
// - application/json
// parameters:
// - name: body
//   in: body
//   description: the register obj
//   schema:
//    $ref: '#/definitions/RegisterRequest'
// responses:
//  '201':
//    description: successful operation
//    schema:
//     $ref: '#/definitions/Response'
//  '400':
//    description: invalid request
//    schema:
//     $ref: '#/definitions/Response'
//  '500':
//    description: internal server error
//    schema:
//     $ref: '#/definitions/Response'
func (c *AppController) Register(ctx *gin.Context) {
	request := new(models.RegisterRequest)
	if err := ctx.BindJSON(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", err))
		return
	}

	found, err := c.appRepo.GetUserByEmail(ctx, request.Email)

	if err == nil && found != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, getErrorResponse(fmt.Sprintf("user with email [%s] exists", request.Email), err))
		return
	}

	now := time.Now().Unix()
	user := models.User{
		Name:               request.Name,
		Email:              request.Email,
		Password:           request.Password,
		DateUpdated:        now,
		DateCreated:        now,
		FailedLoginAttempt: 0,
		LastLogin:          0,
	}

	if len(user.Email) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", errors.New("field email is required")))
		return
	}

	if len(user.Password) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", errors.New("field password is required")))
		return
	}

	if len(user.Name) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", errors.New("field name is required")))
		return
	}

	bytesPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("failed to register user", err))
		return
	}

	user.Password = string(bytesPassword)

	_, err = c.appRepo.AddUser(ctx, user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("failed to register user", err))
		return
	}

	ctx.JSON(http.StatusCreated, &models.Response{
		Message: "successfully registered user",
	})
}
