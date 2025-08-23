// Package storage internal/storage/database.go
package storage

import (
	"github.com/mohar9h/goauth/internal/utils"
	"time"

	"github.com/mohar9h/goauth/internal/entity"
	"gorm.io/gorm"
)

type gormDriver struct {
	db *gorm.DB
}

func NewGormDriver(db *gorm.DB) Driver {
	return &gormDriver{db: db}
}

func (g *gormDriver) StoreToken(t *entity.PersonalAccessToken) error {
	return g.db.Create(t).Error
}

func (g *gormDriver) FindByID(id int64) (*entity.PersonalAccessToken, error) {
	var t entity.PersonalAccessToken
	if err := g.db.First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
		return nil, utils.ErrTokenExpired
	}
	return &t, nil
}

func (g *gormDriver) FindByHash(hash string) (*entity.PersonalAccessToken, error) {
	var t entity.PersonalAccessToken

	if err := g.db.First(&t, "token = ?", hash).Error; err != nil {
		return nil, err
	}
	if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
		return nil, utils.ErrTokenExpired
	}
	return &t, nil
}

func (g *gormDriver) RevokeToken(hash string) error {
	return g.db.Delete(&entity.PersonalAccessToken{}, "token = ?", hash).Error
}

func (g *gormDriver) TouchLastUsed(id int64) error {
	return g.db.Model(&entity.PersonalAccessToken{}).
		Where("id = ?", id).
		Update("last_used_at", time.Now()).
		Error
}
