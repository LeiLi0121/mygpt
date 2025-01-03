package repository

import (
	"mygpt-back/user-service/internal/model"

	"gorm.io/gorm"
)

// Dependency Injection 依赖注入设计模式
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建 UserRepository 实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser 创建用户
func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 未找到记录，返回 nil
		}
		return nil, err // 返回其他数据库错误
	}
	return &user, nil // 找到记录
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
