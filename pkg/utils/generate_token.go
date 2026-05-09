package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	StaffID    uuid.UUID `json:"staff_id"`
	HospitalID uuid.UUID `json:"hospital_id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(staffID, hospitalID uuid.UUID, username, role string) (string, time.Time, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", time.Time{}, fmt.Errorf("JWT_SECRET is not set")
	}

	expiresIn := 24 * time.Hour
	if d := os.Getenv("JWT_EXPIRES_IN"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			expiresIn = parsed
		}
	}

	expiresAt := time.Now().Add(expiresIn)

	claims := JWTClaims{
		StaffID:    staffID,
		HospitalID: hospitalID,
		Username:   username,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   staffID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, expiresAt, nil
}

func ParseToken(tokenStr string) (*JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
