package auth

import (
	"errors"
	"fmt"
	"log"
	"github.com/box1bs/TelegraphicVault/pkg/config"
	"github.com/box1bs/TelegraphicVault/pkg/database"
	"github.com/box1bs/TelegraphicVault/pkg/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrAuthFailed = errors.New("invalid username or password")

type Claims struct {
    UserID    uuid.UUID `json:"user_id"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    jwt.RegisteredClaims
}

type AuthService struct {
	config *config.AuthConfig
	store  storage.JWTUserStorage
}

type tokenPair struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    string    `json:"expires_at"`
}

func NewAuthService(config *config.AuthConfig, store storage.JWTUserStorage) *AuthService {
    return &AuthService{
        config: config,
        store:  store,
    }
}

func (s *AuthService) Login(username, password string) (*tokenPair, error) {
    user, err := s.store.FindByUsername(username)
    if err != nil {
        return nil, ErrAuthFailed
    }

    if err := bcrypt.CompareHashAndPassword(
        []byte(user.Password),
        []byte(password),
    ); err != nil {
        return nil, ErrAuthFailed
    }

    if err := s.store.LastLoginUpdate(user); err != nil {
        log.Printf("last_login update failed for %s: %v\n", user.ID, err)
    }

    return s.generateTokenPair(user)
}

func (s *AuthService) Register(u *model.User) (*tokenPair, error) {
    if err := s.store.LastLoginUpdate(u); err != nil {
        log.Printf("last_login update failed for %s: %v\n", u.ID, err)
    }
    return s.generateTokenPair(u)
}

func (s *AuthService) generateTokenPair(user *model.User) (*tokenPair, error) {
    accessToken, expiresAt, err := s.generateAccessToken(user)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.generateRefreshToken(user)
    if err != nil {
        return nil, err
    }

    return &tokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt.String(),
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
            Subject:   user.ID.String(),
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
        Subject:   user.ID.String(),
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

        c.Set("user_id", claims.UserID.String())
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

func (s *AuthService) RefreshTokens(refreshToken string) (*tokenPair, error) {
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

    return s.generateTokenPair(user)
}

func extractToken(c *gin.Context) string {
    bearerToken := c.GetHeader("Authorization")
    if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
        return bearerToken[7:]
    }
    return ""
}