package v1

import (
	"net/http"

	"URL-Shortner/models"

	"URL-Shortner/business"
	"URL-Shortner/constants"
	"URL-Shortner/log"

	"github.com/gin-gonic/gin"
)

// HandleRegisterUser godoc
// // @Summary Register a new user
// // @Description Register a new user with username and password
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param user body models.User true "User registration request"
// // @Success 200 {object} models.User "User object"
// // @Failure 400 {object} constants.Error "Bad request"
// // @Failure 500 {object} constants.Error "Internal server error"
// // @Router /v1/register [post]
func HandleRegisterUser(ctx *gin.Context) {
	var userReq models.User
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrBindJSONFailed.SetErr(err))
		return
	}
	// Validate the user data
	if err := userReq.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// Register the user
	user, err := business.RegisterUser(userReq)
	if err != nil {
		log.Sugar.Errorf("Failed to register user: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
