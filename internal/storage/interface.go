package storage

import (
	"github.com/mohar9h/goauth/internal/entity"
)

type Driver interface {
	FindByID(id int64) (*entity.Token, error)
	FindByHash(hash string) (*entity.Token, error)
	RevokeToken(hash string) error
	TouchLastUsed(id int64) error
	StoreToken(t *entity.Token) error
}
