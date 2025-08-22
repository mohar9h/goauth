// Package entity internal/entity/session.go
package entity

import "time"

type Session struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	UserId         int64  `gorm:"index"`
	TokenID        string `gorm:"index"`    // hashed token
	ClientInfo     string `gorm:"size:255"` // browser/device
	IP             string `gorm:"size:45"`  // IPv6-safe
	CreatedAt      time.Time
	LastUsedAt     *time.Time
	IdleTimeout    time.Duration // example 30 minutes
	AbsoluteExpiry time.Time     // example 30 days
}
