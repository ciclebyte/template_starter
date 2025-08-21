package libPassword

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// Params 密码哈希参数
type Params struct {
	Memory      uint32 // 内存使用量 (KB)
	Iterations  uint32 // 迭代次数
	Parallelism uint8  // 并行度
	SaltLength  uint32 // 盐值长度
	KeyLength   uint32 // 哈希长度
}

// 默认参数配置
var defaultParams = &Params{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword 对密码进行哈希
func HashPassword(password string) (string, error) {
	return HashPasswordWithParams(password, defaultParams)
}

// HashPasswordWithParams 使用指定参数对密码进行哈希
func HashPasswordWithParams(password string, p *Params) (string, error) {
	// 生成随机盐值
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		g.Log().Error(context.Background(), "failed to generate salt:", err)
		return "", err
	}

	// 使用Argon2id算法生成哈希
	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// 编码为base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式化为存储格式: $argon2id$v=19$m=65536,t=3,p=2$saltBase64$hashBase64
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析存储的哈希
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// 使用相同参数对输入密码进行哈希
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// 使用constant-time比较防止时序攻击
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// generateRandomBytes 生成随机字节
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash 解析存储的哈希字符串
func decodeHash(encodedHash string) (p *Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// IsValidPassword 检查密码强度
func IsValidPassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "密码长度至少6位"
	}
	
	if len(password) > 64 {
		return false, "密码长度不能超过64位"
	}

	// 检查是否包含基本字符类型
	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	// 至少包含两种字符类型
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}

	if count < 2 {
		return false, "密码必须包含至少两种字符类型（大写字母、小写字母、数字）"
	}

	return true, ""
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) (string, error) {
	if length < 6 {
		length = 6
	}
	if length > 64 {
		length = 64
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b), nil
}

// CheckPasswordHistory 检查密码是否在历史记录中（简单实现）
func CheckPasswordHistory(password string, historyHashes []string) (bool, error) {
	for _, historyHash := range historyHashes {
		match, err := VerifyPassword(password, historyHash)
		if err != nil {
			continue // 忽略无效的历史哈希
		}
		if match {
			return true, nil // 密码在历史记录中
		}
	}
	return false, nil
}

// EstimateHashingTime 估算哈希时间（用于性能调优）
func EstimateHashingTime(password string, p *Params) (string, error) {
	start := g.TimestampMilli()
	_, err := HashPasswordWithParams(password, p)
	if err != nil {
		return "", err
	}
	duration := g.TimestampMilli() - start
	return strconv.FormatInt(duration, 10) + "ms", nil
}