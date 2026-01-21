package users

import (
	"alloy/internal/shared/database/models"
	"context"

	"fmt"

	"strings"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	FindAllUsersWithEmails(ctx context.Context, emails []string) ([]models.User, error)
	GetUserByEmployeeNumber(ctx context.Context, employeeNumber string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) FindAllUsersWithEmails(
	ctx context.Context,
	emails []string,
) ([]models.User, error) {

	if len(emails) == 0 {
		return nil, errors.New("emails list cannot be empty")
	}

	// Normalize (important!)
	normalized := make([]string, 0, len(emails))
	seen := make(map[string]struct{})

	for _, e := range emails {
		e = strings.TrimSpace(strings.ToLower(e))
		if e == "" {
			return nil, errors.New("email cannot be empty")
		}
		if _, ok := seen[e]; !ok {
			seen[e] = struct{}{}
			normalized = append(normalized, e)
		}
	}

	var users []models.User

	err := r.db.
		WithContext(ctx).
		Where("email IN ?", normalized).
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}

	// Build lookup of returned emails
	found := make(map[string]struct{}, len(users))
	for _, u := range users {
		found[strings.ToLower(u.Email)] = struct{}{}
	}

	// Detect missing emails
	var missing []string
	for _, e := range normalized {
		if _, ok := found[e]; !ok {
			missing = append(missing, e)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf(
			"%w: %v",
			errors.New("one or more emails do not exist"),
			missing,
		)
	}

	return users, nil
}

func (r *userRepository) GetUserByEmployeeNumber(ctx context.Context, employeeNumber string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("employee_number = ?", employeeNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
