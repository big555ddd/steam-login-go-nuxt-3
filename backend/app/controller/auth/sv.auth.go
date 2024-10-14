package auth

import (
	"app/app/helper"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GenerateToken(ctx context.Context, authType string, user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth_type": authType,
		"id":        user.ID,
		"data": map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"nickname":  user.Nickname,
			"status":    user.Status,
		},
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(7 * viper.GetDuration("TOKEN_DURATION_USER")).Unix(),
	})

	secret := []byte(viper.GetString("TOKEN_SECRET_USER"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.Infof("[error]: %v", err)
		return "", err
	}
	return tokenString, nil
}

func (s *Service) GenerateSteamToken(ctx context.Context, player *helper.PlayerProfile) (string, error) {
	claims := jwt.MapClaims{
		"steam_id":     player.SteamID,
		"personaname":  player.PersonaName,
		"avatarsmall":  player.AvatarSmall,
		"avatarmedium": player.AvatarMedium,
		"avatarfull":   player.AvatarFull,
		"nbf":          time.Now().Unix(),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(viper.GetString("JWT_SECRET")) // อ่าน JWT_SECRET จาก config
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Printf("[error]: %v", err)
		return "", err
	}
	return tokenString, nil
}

func (s *Service) Login(ctx context.Context, req model.User) (*model.User, error) {
	storedUser, err := s.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Log the stored hashed password and the incoming plain text password
	log.Printf("Stored password hash: %s", storedUser.Password)
	log.Printf("Incoming password: %s", req.Password)

	// Check if the provided password matches the stored password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return storedUser, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	m := model.User{}
	logger.Infof("Searching for user with username: %s", username)
	if err := s.db.NewSelect().Model(&m).
		Where("username = ?", username).Scan(ctx); err != nil {
		logger.Infof("[error]: %v", err)
		return nil, err
	}
	return &m, nil
}

func (s *Service) GetUserBySteamID(ctx context.Context, steam_id string) (*model.SteamUser, error) {
	m := model.SteamUser{}
	logger.Infof("Searching for user with steam_id: %s", steam_id)
	if err := s.db.NewSelect().Model(&m).
		Where("steam_id = ?", steam_id).Scan(ctx); err != nil {
		logger.Infof("[error]: %v", err)
		return nil, err
	}
	return &m, nil
}

func (s *Service) CreateSteamUser(ctx context.Context, steam *helper.PlayerProfile) (model.SteamInfo, error) {
	m := model.SteamInfo{
		SteamID:      steam.SteamID,
		PersonaName:  steam.PersonaName,
		AvatarSmall:  steam.AvatarSmall,
		AvatarMedium: steam.AvatarMedium,
		AvatarFull:   steam.AvatarFull,
	}

	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		logger.Infof("[error]: %v", err)
		return m, err
	}

	newSteamUser := model.SteamUser{
		SteamID:   m.SteamID,
		Points:    0,
		VipPoints: 0,
	}

	_, err = s.db.NewInsert().Model(&newSteamUser).Exec(ctx)
	if err != nil {
		logger.Infof("[error]: %v", err)
		return m, err
	}

	return m, nil
}

func (s *Service) GetUserDetailByToken(ctx context.Context, tokenString string) (response.GetUserDetail, error) {
	var userDetail response.GetUserDetail

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(viper.GetString("TOKEN_SECRET_USER")), nil
	})
	if err != nil {
		return userDetail, err
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userDetail, errors.New("invalid token")
	}

	data, ok := claims["data"].(map[string]interface{})
	if !ok {
		return userDetail, errors.New("invalid token payload: data field is missing")
	}

	userIDFloat, ok := data["id"].(float64)
	if !ok {
		return userDetail, errors.New("invalid token payload: id field is missing or not a number")
	}
	userID := int(userIDFloat)

	// Define the query and execute it using Bun
	var user model.User
	err = s.db.NewSelect().
		Model(&user).
		Column("id", "username", "firstname", "lastname", "nickname", "email", "role_id", "points").
		Where("id = ?", userID).
		Scan(ctx)
	if err != nil {
		logger.Infof("[error]: Failed to fetch user: %v", err)
		return userDetail, err
	}

	// userDetail = response.GetUserDetail{
	// 	ID:        user.ID,
	// 	Username:  user.Username,
	// 	Firstname: user.Firstname,
	// 	Lastname:  user.Lastname,
	// 	Nickname:  user.Nickname,
	// 	Email:     user.Email,
	// 	RoleID:    user.RoleID,
	// 	Point:     user.Points,
	// }

	return userDetail, nil
}

func (s *Service) LoginWithSteam(ctx context.Context, steamLogin request.SteamLogin) (*model.User, error) {
	// Logic สำหรับตรวจสอบ SteamID และดึงข้อมูลผู้ใช้
	// user, err := s.db.FindUserBySteamID(ctx, steamLogin.SteamID)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
