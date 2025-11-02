// Package storage internal/storage/factory.go
package storage

import (
	"errors"
	"gorm.io/gorm"
)

type Type string

const (
	Memory Type = "memory"
	Gorm   Type = "gorm"
	Redis  Type = "redis"
)

// Config configuration for storage drivers
type Config struct {
	Type   Type
	GormDB *gorm.DB
	// RedisClient *redis.Client // اگر Redis اضافه کنید
}

// NewStorage creates a new storage driver based on configuration
func NewStorage(config Config) (Driver, error) {
	switch config.Type {
	case Memory:
		return NewMemoryDriver(), nil
	case Gorm:
		if config.GormDB == nil {
			return nil, errors.New("gorm DB instance is required")
		}
		return NewGormDriver(config.GormDB), nil
	case Redis:
		// if config.RedisClient == nil {
		//     return nil, errors.New("redis client is required")
		// }
		// return NewRedisDriver(config.RedisClient), nil
		return nil, errors.New("redis driver not implemented yet")
	default:
		return nil, errors.New("unknown storage type")
	}
}
