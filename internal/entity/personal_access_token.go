// Package entity internal/entity/personal_access_token.go
package entity

import "time"

// PersonalAccessToken APIToken defines the persistent structure stored in SQL/Redis.
type PersonalAccessToken struct {
	ID         int64      `gorm:"primaryKey;autoIncrement"`
	UserId     int64      `gorm:"index"`
	Token      string     `gorm:"index;size:100"`
	Name       *string    `gorm:"size:100"`
	Abilities  string     `gorm:"type:text"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	ExpiresAt  *time.Time `gorm:"index"`
	LastUsedAt *time.Time
}

func (PersonalAccessToken) TableName() string { return "personal_access_tokens" }
