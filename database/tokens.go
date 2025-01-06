package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"something/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (p *Postgres) SaveRefreshToken(userID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	token := model.UserToken{
		UserID:     userID,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		Created_at: time.Now(),
	}

	return p.db.Create(&token).Error
}

func (p *Postgres) ValidateRefreshToken(userID uuid.UUID, refreshToken string) (bool, error) {
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	var token model.UserToken
	err := p.db.Where("user_id = ? AND token_hash = ? AND is_revoked = ? AND expires_at > ?", userID, tokenHash, false, time.Now()).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *Postgres) RevokeRefreshToken(userID uuid.UUID, refreshToken string) error {
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	return p.db.Model(&model.UserToken{}).Where("user_id = ? AND token_hash = ?", userID, tokenHash).Update("is_revoked", true).Error
}