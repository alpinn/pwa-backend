package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pwa-backend/internal/models"
	"pwa-backend/internal/repositories"
)

type AuthHandler struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
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

func NewAuthHandler(userRepo *repositories.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
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
func (h *AuthHandler) Login(c *gin.Context) {
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

	now := time.Now()
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = "powersync-key" 
	
	token.Claims = jwt.MapClaims{
		"sub":      user.ID,                       
		"aud":      "https://6940ebf14011d65924582a54.powersync.journeyapps.com", 
		"iss":      "pwa-backend",                 
		"jti":      fmt.Sprintf("%s-%d", user.ID, now.Unix()),
		"nbf":      now.Unix(),                    
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"iat":      now.Unix(),                    
		"exp":      now.Add(24 * time.Hour).Unix(),
	}

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenString,
		User:  *user,
	})
}

// JWKS godoc
// @Summary Get JSON Web Key Set
// @Description Get JWKS for PowerSync integration
// @Tags auth
// @Produce json
// @Success 200 {object} JWKSResponse
// @Router /.well-known/jwks.json [get]
func (h *AuthHandler) JWKS(c *gin.Context) {
	keyBytes := []byte(h.jwtSecret)
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
func (h *AuthHandler) PowerSyncAuth(c *gin.Context) {
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
		return []byte(h.jwtSecret), nil
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