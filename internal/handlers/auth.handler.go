package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pwa-backend/internal/config"
	"pwa-backend/internal/models"
	"pwa-backend/internal/repositories"
)

type AuthHandlerStruct struct {
	userRepo          *repositories.UserRepository
	sessionRepo       *repositories.UserSessionRepository
    jwtConfig            *config.JWTConfig
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	K   string `json:"k"`
}

func AuthHandler(
	userRepo *repositories.UserRepository,
	sessionRepo *repositories.UserSessionRepository,
    jwtConfig *config.JWTConfig,
) *AuthHandlerStruct {
	return &AuthHandlerStruct{
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		jwtConfig:            jwtConfig,
		accessTokenDuration:  10 * time.Minute,
		refreshTokenDuration: 7 * 24 * time.Hour,
	}
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandlerStruct) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if user already exists
    existingUser, err := h.userRepo.GetByUsername(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    if existingUser != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := &models.User{
        ID:        uuid.New().String(),
        Username:  req.Username,
        Password:  string(hashedPassword),
        Name:      req.Name,
        Role:      req.Role,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if user.Role == "" {
        user.Role = "staff"
    }

    if err := h.userRepo.Create(user); err != nil {
        log.Printf("Failed to create user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Create session
    sessionID := uuid.New().String()
    refreshToken := uuid.New().String()
    now := time.Now()

    session := &models.UserSession{
        ID:           sessionID,
        UserID:       user.ID,
        RefreshToken: refreshToken,
        ExpiresAt:    now.Add(h.refreshTokenDuration),
        CreatedAt:    now,
    }

    if err := h.sessionRepo.Create(session); err != nil {
        log.Printf("Failed to create session: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
        return
    }

    // Generate access token
    accessToken := jwt.New(jwt.SigningMethodHS256)
    accessToken.Header["kid"] = "powersync-key"

    accessToken.Claims = jwt.MapClaims{
        "sub":      user.ID,
        "aud":      h.jwtConfig.Audience,
        "iss":      h.jwtConfig.Issuer,
        "jti":      fmt.Sprintf("%s-%d", user.ID, now.Unix()),
        "nbf":      now.Unix(),
        "user_id":  user.ID,
        "username": user.Username,
        "role":     user.Role,
        "sid":      sessionID,
        "iat":      now.Unix(),
        "exp":      now.Add(h.accessTokenDuration).Unix(),
    }

    accessTokenString, err := accessToken.SignedString([]byte(h.jwtConfig.Secret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Return user without password
    userResponse := *user
    userResponse.Password = ""

    c.JSON(http.StatusCreated, models.TokenResponse{
        AccessToken:  accessTokenString,
        RefreshToken: refreshToken,
        User:         userResponse,
    })
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandlerStruct) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create session
	sessionID := uuid.New().String()
	refreshToken := uuid.New().String()
	now := time.Now()
	
	session := &models.UserSession{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(h.refreshTokenDuration),
		CreatedAt:    now,
	}
	
	if err := h.sessionRepo.Create(session); err != nil {
		log.Printf("Failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Generate access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessToken.Header["kid"] = "powersync-key"
	
	accessToken.Claims = jwt.MapClaims{
		"sub":      user.ID,
		"aud":      h.jwtConfig.Audience,
        "iss":      h.jwtConfig.Issuer,
		"jti":      fmt.Sprintf("%s-%d", user.ID, now.Unix()),
		"nbf":      now.Unix(),
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"sid":      sessionID,
		"iat":      now.Unix(),
		"exp":      now.Add(h.accessTokenDuration).Unix(),
	}

    accessTokenString, err := accessToken.SignedString([]byte(h.jwtConfig.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Use refresh token to get new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandlerStruct) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get session by refresh token
	session, err := h.sessionRepo.GetByRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if session == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	user, err := h.userRepo.GetByID(session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	now := time.Now()
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = "powersync-key"

	token.Claims = jwt.MapClaims{
		"sub":      user.ID,
		"aud":      h.jwtConfig.Audience,
        "iss":      h.jwtConfig.Issuer,
		"jti":      fmt.Sprintf("%s-%d", user.ID, now.Unix()),
		"nbf":      now.Unix(),
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"sid":      session.ID,
		"iat":      now.Unix(),
		"exp":      now.Add(h.accessTokenDuration).Unix(),
	}

	accessTokenString, err := token.SignedString([]byte(h.jwtConfig.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: req.RefreshToken,
		User:         *user,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Revoke user session
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandlerStruct) Logout(c *gin.Context) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
		return
	}

	if err := h.sessionRepo.RevokeByID(sessionID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// JWKS godoc
// @Summary Get JSON Web Key Set
// @Description Get JWKS for PowerSync integration
// @Tags auth
// @Produce json
// @Success 200 {object} JWKSResponse
// @Router /.well-known/jwks.json [get]
func (h *AuthHandlerStruct) JWKS(c *gin.Context) {
	keyBytes := []byte(h.jwtConfig.Secret)
	encoded := base64.RawURLEncoding.EncodeToString(keyBytes)

	log.Printf("JWKS Endpoint Called - KID: %s, Secret (base64): %s", "powersync-key", encoded)

	jwks := JWKSResponse{
		Keys: []JWK{
			{
				Kty: "oct",           
				Use: "sig",           
				Kid: "powersync-key", 
				Alg: "HS256",         
				K:   encoded,         
			},
		},
	}

	c.JSON(http.StatusOK, jwks)
}

// PowerSyncAuth godoc
// @Summary Get PowerSync Auth endpoint
// @Description Endpoint specifically for PowerSync authentication
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/powersync [get]
func (h *AuthHandlerStruct) PowerSyncAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.jwtConfig.Secret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{
			"user_id":  claims["user_id"],
			"username": claims["username"],
			"role":     claims["role"],
			"valid":    true,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	}
}

// Me godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandlerStruct) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, err := h.userRepo.GetByID(userID.(string))
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}