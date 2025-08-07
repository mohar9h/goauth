package auth

import (
	"github.com/mohar9h/goauth/config"
	"gorm.io/gorm"
)

type TokenOptions struct {
	UserId    int64
	Name      *string
	Abilities []string
	Config    *config.Config
	DB        *gorm.DB // Required for GORM storage
}
