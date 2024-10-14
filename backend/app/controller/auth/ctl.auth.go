package auth

import (
	"app/app/helper"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Login(c *gin.Context) {
	var loginUser request.LoginUser
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a context
	ctx := context.Background()

	// Convert loginUser to model.User
	user := model.User{
		Username: loginUser.Username,
		Password: loginUser.Password,
	}

	loggedInUser, err := ctl.Service.Login(ctx, user)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// Generate a token for the logged-in user
	token, err := ctl.Service.GenerateToken(ctx, loggedInUser.Username, loggedInUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response.Success(c, token)
}

func (ctl *Controller) GetUserDetailByToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userDetail, err := ctl.Service.GetUserDetailByToken(c.Request.Context(), tokenString)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, userDetail)
}

// LoginWithSteam ใช้ Service เพื่อสร้าง JWT Token
func (ctl *Controller) LoginWithSteam(c *gin.Context) {
	steamID := c.Param("id")
	if steamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Steam ID"})
		return
	}
	summaries, err := helper.GetSteamProfile(steamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := ctl.Service.GenerateSteamToken(c.Request.Context(), summaries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = ctl.Service.GetUserBySteamID(c.Request.Context(), summaries.SteamID)
	if err != nil {
		// Create a new user
		_, err = ctl.Service.CreateSteamUser(c.Request.Context(), summaries)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"profile": summaries,
	})
}
