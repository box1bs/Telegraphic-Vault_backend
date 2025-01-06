package auth

import (
	"errors"
	"fmt"
	"something/config"
	"something/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
    UserID    uuid.UUID `json:"user_id"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    jwt.RegisteredClaims
}

type AuthService struct {
	config config.AuthConfig
	store  config.UserStore
}

type TokenPair struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}

func NewAuthService(config config.AuthConfig, store config.UserStore) *AuthService {
    return &AuthService{
        config: config,
        store:  store,
    }
}

func (s *AuthService) Login(username, password string) (*TokenPair, error) {
    user, err := s.store.FindByUsername(username)
    if err != nil {
        return nil, errors.New("invalid username or password")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid username or password")
    }

    return s.GenerateTokenPair(user)
}

func (s *AuthService) GenerateTokenPair(user *model.User) (*TokenPair, error) {
    accessToken, expiresAt, err := s.generateAccessToken(user)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.generateRefreshToken(user)
    if err != nil {
        return nil, err
    }

    if err := s.store.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(s.config.RefreshTokenTTL)); err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, time.Time, error) {
    expiresAt := time.Now().Add(s.config.AccessTokenTTL)
    
    claims := Claims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   fmt.Sprintf("%d", user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    tokenString, err := token.SignedString([]byte(s.config.AccessTokenSecret))
    if err != nil {
        return "", time.Time{}, err
    }

    return tokenString, expiresAt, nil
}

func (s *AuthService) generateRefreshToken(user *model.User) (string, error) {
    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        Subject:   fmt.Sprintf("%d", user.ID),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.RefreshTokenSecret))
}

func (s *AuthService) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        claims, err := s.validateAccessToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}

func (s *AuthService) validateAccessToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.config.AccessTokenSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

func (s *AuthService) RefreshTokens(refreshToken string) (*TokenPair, error) {
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.config.RefreshTokenSecret), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid refresh token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }

    userID, err := uuid.Parse(claims["sub"].(string))
    if err != nil {
        return nil, err
    }
    user, err := s.store.FindByID(userID)
    if err != nil {
        return nil, err
    }

    if err := s.store.RevokeRefreshToken(user.ID, refreshToken); err != nil {
        return nil, err
    }

    return s.GenerateTokenPair(user)
}

func extractToken(c *gin.Context) string {
    bearerToken := c.GetHeader("Authorization")
    if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
        return bearerToken[7:]
    }
    return ""
}