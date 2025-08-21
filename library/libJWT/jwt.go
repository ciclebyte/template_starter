package libJWT

import (
	"context"
	"errors"
	"time"

	"github.com/ciclebyte/template_starter/library/libConfig"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenNotValidYet = errors.New("token used before valid")
)

// JWTManager JWT管理器
type JWTManager struct {
	SecretKey     []byte
	AccessExpire  time.Duration // 访问令牌过期时间 (2小时)
	RefreshExpire time.Duration // 刷新令牌过期时间 (7天)
	Issuer        string
}

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Claims JWT载荷
type Claims struct {
	UserID      int64    `json:"user_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"` // access 或 refresh
	jwt.RegisteredClaims
}

var jwtManager *JWTManager

// Init 初始化JWT管理器
func Init() {
	ctx := context.Background()
	
	secretKey := libConfig.GetString(ctx, "jwt.secret_key", "your-256-bit-secret-key-change-in-production!")
	accessExpire := libConfig.GetInt(ctx, "jwt.access_expire", 7200) // 2小时
	refreshExpire := libConfig.GetInt(ctx, "jwt.refresh_expire", 604800) // 7天
	issuer := libConfig.GetString(ctx, "jwt.issuer", "template_starter")

	jwtManager = &JWTManager{
		SecretKey:     []byte(secretKey),
		AccessExpire:  time.Duration(accessExpire) * time.Second,
		RefreshExpire: time.Duration(refreshExpire) * time.Second,
		Issuer:        issuer,
	}

	g.Log().Info(context.Background(), "JWT Manager initialized")
}

// GetManager 获取JWT管理器实例
func GetManager() *JWTManager {
	if jwtManager == nil {
		Init()
	}
	return jwtManager
}

// GenerateTokens 生成访问令牌和刷新令牌
func (j *JWTManager) GenerateTokens(userID int64, username, email string, roles, permissions []string) (*TokenInfo, error) {
	now := time.Now()
	
	// 生成访问令牌
	accessClaims := &Claims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.AccessExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.Issuer,
			Subject:   username,
			ID:        guid.S(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.SecretKey)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshClaims := &Claims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.Issuer,
			Subject:   username,
			ID:        guid.S(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.SecretKey)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(j.AccessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken 验证令牌
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保token方法是HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.SecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshAccessToken 刷新访问令牌
func (j *JWTManager) RefreshAccessToken(refreshTokenString string, roles, permissions []string) (*TokenInfo, error) {
	// 验证刷新令牌
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// 确保这是一个刷新令牌
	if claims.TokenType != "refresh" {
		return nil, ErrTokenInvalid
	}

	// 生成新的令牌对
	return j.GenerateTokens(claims.UserID, claims.Username, claims.Email, roles, permissions)
}

// ExtractTokenFromHeader 从Authorization头中提取token
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}