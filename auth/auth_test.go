package auth_test

import (
	"fmt"
	"github.com/mohar9h/goauth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

type PostgresConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func OpenPostgres(cfg PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // تنظیم لاگ (یا .Silent / .Info)
	})
	if err != nil {
		return nil, err
	}

	// اتصال پایه‌ای برای اعمال تنظیمات connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Minute)

	return db, nil
}

//func TestCreateRandomToken(t *testing.T) {
//	db, err := OpenPostgres(PostgresConfig{
//		Host:            "localhost",
//		Port:            5432,
//		Username:        "postgres",
//		Password:        "Mohammad1367",
//		Database:        "pecalets",
//		SSLMode:         "disable",
//		MaxIdleConns:    15,
//		MaxOpenConns:    100,
//		ConnMaxLifetime: 5,
//	})
//	if err != nil {
//		t.Fatalf("failed to connect to PostgreSQL: %v", err)
//	}
//
//	// مهاجرت جدول مورد نیاز برای توکن‌ها
//	if err := db.AutoMigrate(&goauth.PersonalAccessToken{}); err != nil {
//		t.Fatalf("failed to migrate token table: %v", err)
//	}
//
//	// فراخوانی تابع تولید توکن
//	result, err := goauth.CreateToken(&goauth.TokenOptions{
//		UserId:    1,
//		Name:      nil,
//		Abilities: []string{"*"},
//		DB:        db,
//	})
//
//	if err != nil {
//		t.Fatalf("CreateToken failed: %v", err)
//	}
//
//	// خروجی توکن لاگ شود
//	t.Logf("Generated Token: %s", result)
//}

func TestValidateToken(t *testing.T) {

	db, err := OpenPostgres(PostgresConfig{
		Host:            "localhost",
		Port:            5432,
		Username:        "postgres",
		Password:        "Mohammad1367",
		Database:        "pecalets",
		SSLMode:         "disable",
		MaxIdleConns:    15,
		MaxOpenConns:    100,
		ConnMaxLifetime: 5,
	})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	tokenStr := "Bearer 7|93cfb2035a4ef58f290f60e7c1cb3401d7dc786c189d8754"

	cfg := goauth.SetupGorm(db)

	validateToken, err := goauth.ValidateToken(tokenStr, cfg)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	t.Logf("Valid Token: %+v", validateToken)
}

//func TestRevokeToken(t *testing.T) {
//	db, err := OpenPostgres(PostgresConfig{
//		Host:            "localhost",
//		Port:            5432,
//		Username:        "postgres",
//		Password:        "Mohammad1367",
//		Database:        "pecalets",
//		SSLMode:         "disable",
//		MaxIdleConns:    15,
//		MaxOpenConns:    100,
//		ConnMaxLifetime: 5,
//	})
//	if err != nil {
//		t.Fatalf("failed to connect to PostgreSQL: %v", err)
//	}
//
//	tokenStr, err := goauth.CreateToken(&goauth.TokenOptions{
//		UserId:    1,
//		Abilities: []string{"*"},
//		DB:        db, // فرض: اینجا DB از قبل آماده است
//	})
//	if err != nil {
//		t.Fatalf("CreateToken failed: %v", err)
//	}
//
//	cfg := goauth.SetupGorm(db)
//	// باطل‌کردن توکن
//	err = goauth.RevokeToken(tokenStr, cfg)
//	if err != nil {
//		t.Fatalf("RevokeToken failed: %v", err)
//	}
//
//	// بررسی اعتبارسنجی بعد از ابطال
//	_, err = goauth.ValidateToken(tokenStr, nil)
//	if err == nil {
//		t.Fatalf("Token should be revoked but is still valid")
//	}
//
//	t.Logf("Token successfully revoked: %v", tokenStr)
//}
