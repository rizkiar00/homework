package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkiar00/homework/pkg/config"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	secret    []byte
	expiresIn time.Duration
}

func New(cfg config.Config) *Service {
	return &Service{
		secret:    []byte(cfg.JWT.Secret),
		expiresIn: cfg.JWT.ExpiresIn(),
	}
}

func (s *Service) Generate(userID, username, role string) (string, int, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, int(s.expiresIn.Seconds()), nil
}

func (s *Service) Parse(value string) (Claims, error) {
	claims := Claims{}
	parsed, err := jwt.ParseWithClaims(value, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !parsed.Valid {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
