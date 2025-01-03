package service

import (
	"context"
	"errors"
	"mygpt-back/user-service/internal/model"
	"mygpt-back/user-service/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  *repository.UserRepository
	redis *redis.Client
}

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// NewUserService 创建 UserService 实例
func NewUserService(repo *repository.UserRepository, redis *redis.Client) *UserService {
	return &UserService{repo: repo, redis: redis}
}

// RegisterUser 用户注册
func (s *UserService) RegisterUser(user *model.User) error {
	// 检查用户名是否存在
	existingUser, _ := s.repo.GetUserByUsername(user.Username)
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已注册
	existingEmailUser, _ := s.repo.GetUserByEmail(user.Email)
	if existingEmailUser != nil {
		return errors.New("邮箱已被注册")
	}

	// 密码加密
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// 保存用户
	return s.repo.CreateUser(user)
}

// LoginUser 用户登录
func (s *UserService) LoginUser(username, password string) (string, error) {
	// 根据用户名获取用户
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("用户名不存在")
	}

	if user == nil {
		return "", errors.New("用户不存在")
	}

	// 验证密码
	if !CheckPassword(password, user.Password) {
		return "", errors.New("密码错误")
	}

	// 生成 JWT Token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	// 保存到 Redis
	err = s.saveTokenToRedis(user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID 根据用户 ID 获取用户信息
func (s *UserService) GetUserByID(userID int) (*model.User, error) {
	// 调用 repository 层获取用户信息
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// LogoutUser 用户登出
func (s *UserService) LogoutUser(token string) error {
	// 删除 Redis 中的 Token
	_, err := s.redis.Del(context.Background(), token).Result()
	return err
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// generateToken 生成 JWT Token
func (s *UserService) generateToken(userID int) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token 24 小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key")) // 替换为实际的密钥
}

// saveTokenToRedis 将 Token 保存到 Redis
func (s *UserService) saveTokenToRedis(userID int, token string) error {
	key := token
	return s.redis.Set(context.Background(), key, userID, 24*time.Hour).Err() // Token 有效期 24 小时
}
