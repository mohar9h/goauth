package storage

import (
	"errors"
	"github.com/mohar9h/goauth/internal/entity"
	"gorm.io/gorm"
	"time"
)

type gormDriver struct {
	db *gorm.DB
}

func NewGormDriver(db *gorm.DB) Driver {
	return &gormDriver{db: db}
}

func (g *gormDriver) StoreToken(t *entity.Token) error {
	return g.db.Create(t).Error
}

func (g *gormDriver) FindByID(id int64) (*entity.Token, error) {

	var t entity.Token
	if err := g.db.First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return &t, nil
}

func (g *gormDriver) FindByHash(hash string) (*entity.Token, error) {
	var t entity.Token

	if err := g.db.First(&t, "token = ?", hash).Error; err != nil {
		return nil, err
	}
	if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return &t, nil
}

func (g *gormDriver) RevokeToken(id string) error {
	return g.db.Delete(&entity.Token{}, "id = ?", id).Error
}

func (g *gormDriver) TouchLastUsed(id int64) error {
	return g.db.Model(&entity.Token{}).
		Where("id = ?", id).
		Update("last_used_at", time.Now()).
		Error
}
