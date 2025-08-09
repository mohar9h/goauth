// Package auth internal/auth/generator.go
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/mohar9h/goauth/config"
	"github.com/mohar9h/goauth/internal/entity"
	"github.com/mohar9h/goauth/internal/utils"
	"hash/crc32"
	"strings"
	"time"
)

type Generator interface {
	Create() (*Result, error)
}

type generator struct {
	opts *TokenOptions
	cfg  *config.Config
}

func NewGenerator(opts *TokenOptions, cfg *config.Config) Generator {
	return &generator{opts: opts, cfg: cfg}
}

func (g *generator) Create() (*Result, error) {
	if g.cfg.Storage == nil {
		return nil, errors.New("no storage backend configured")
	}

	plainText := g.generateTokenString()
	hashed := utils.HashToken(plainText)

	var expireAt *time.Time
	if g.cfg.ExpireAt > 0 {
		t := time.Now().Add(g.cfg.ExpireAt)
		expireAt = &t
	}

	t := &entity.PersonalAccessToken{
		UserId:    g.opts.UserId,
		Name:      g.opts.Name,
		Token:     hashed,
		Abilities: strings.Join(g.opts.Abilities, ","),
		CreatedAt: time.Now(),
		ExpiresAt: expireAt,
	}

	if err := g.cfg.Storage.StoreToken(t); err != nil {
		return nil, err
	}

	return &Result{
		PlainText: fmt.Sprintf("%d|%s", t.ID, plainText),
		TokenID:   hashed,
	}, nil
}

func (g *generator) generateTokenString() string {
	buf := make([]byte, g.cfg.TokenLength)
	if _, err := rand.Read(buf); err != nil {
		panic("token generation failed: " + err.Error())
	}
	raw := hex.EncodeToString(buf)

	crc := crc32.Checksum([]byte(raw), crc32.MakeTable(crc32.Castagnoli))
	return fmt.Sprintf("%s%s%x", g.cfg.TokenPrefix, raw, crc)
}
