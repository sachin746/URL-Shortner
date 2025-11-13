package v1

import (
	"net/http"

	"URL-Shortner/business"
	"URL-Shortner/log"
	"URL-Shortner/models"
	"URL-Shortner/utils"

	"github.com/gin-gonic/gin"
)

// HandleGetUser godoc
// @Summary Get user details
// @Description Retrieve user details by username
// // @Tags user
// // @Accept json
// // // @Produce json
// // // @Param user body models.GetUserRequest true "Get user request"
// // // @Success 200 {object} models.User "User object"
// // // @Failure 400 {object} constants.Error "Bad request"
// // // @Failure 500 {object} constants.Error "Internal server error"
// // @Router /v1/profile [GET]
func HandleGetUser(ctx *gin.Context) {
	var userReq models.GetUserRequest

	username, ok := ctx.Value("username").(string)
	if ok != true {
		log.Sugar.Errorf("Failed to get username from context")
		ctx.JSON(http.StatusBadRequest, "Failed to get username from context")
		return
	}
	userReq.Username = username
	email, ok := ctx.Value("email").(string)
	if !ok {
		log.Sugar.Errorf("Failed to get email from context")
		ctx.JSON(http.StatusBadRequest, "Failed to get email from context")
		return
	}
	userReq.Email = email

	// Validate the user data
	if err := userReq.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// Get the user details
	user, err := business.GetUser(userReq.Username)
	if err != nil {
		log.Sugar.Errorf("Failed to get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(user, "User retrieved successfully"))
}
