package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type SubjectType string

const (
	SubjectUser  SubjectType = "user"
	SubjectAdmin SubjectType = "admin"
)

type Claims struct {
	SubID   string      `json:"sub_id"`   // user_id или admin_id
	SubType SubjectType `json:"sub_type"` // "user" или "admin"
	jwtv5.RegisteredClaims
}

type Manager struct {
	accessSecret []byte
	userTTL      time.Duration
	adminTTL     time.Duration
}

func NewManager(secret string, userTTLMin, adminTTLMin int) *Manager {
	return &Manager{
		accessSecret: []byte(secret),
		userTTL:      time.Duration(userTTLMin) * time.Minute,
		adminTTL:     time.Duration(adminTTLMin) * time.Minute,
	}
}

func (m *Manager) ttlFor(st SubjectType) time.Duration {
	if st == SubjectAdmin {
		return m.adminTTL
	}
	return m.userTTL
}

func (m *Manager) Generate(subID string, subType SubjectType) (string, error) {
	ttl := m.ttlFor(subType)
	claims := Claims{
		SubID:   subID,
		SubType: subType,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   string(subType),
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(m.accessSecret)
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.SubType != SubjectUser && claims.SubType != SubjectAdmin {
		return nil, errors.New("unknown subject type")
	}
	return claims, nil
}

func (m *Manager) UserTTL() time.Duration  { return m.userTTL }
func (m *Manager) AdminTTL() time.Duration { return m.adminTTL }
