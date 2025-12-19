package models

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserStatus enumerates the lifecycle state of an account.
type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

var statusPool = []UserStatus{UserStatusPending, UserStatusActive, UserStatusDisabled}

// User represents a record in the users table.
type User struct {
	ID            uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID          string     `json:"uuid" gorm:"type:char(36);not null;uniqueIndex"`
	Name          string     `json:"name" gorm:"type:varchar(100);not null"`
	Email         string     `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash  string     `json:"password_hash" gorm:"type:varchar(255);not null"`
	AvatarURL     *string    `json:"avatar_url,omitempty" gorm:"type:varchar(512)"`
	Status        UserStatus `json:"status" gorm:"type:enum('pending','active','disabled');default:'pending';not null"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PlainPassword string     `json:"plain_password,omitempty" gorm:"-"`
}

// GenerateFakeUsers returns n fake users backed by gofakeit for quick seed data.
func GenerateFakeUsers(n int) ([]User, error) {
	if n <= 0 {
		return nil, errors.New("count must be greater than zero")
	}

	faker := gofakeit.NewFaker(source.NewCrypto(), false)
	now := time.Now().UTC()

	users := make([]User, n)
	for i := 0; i < n; i++ {
		person := faker.Person()
		name := strings.TrimSpace(fmt.Sprintf("%s %s", person.FirstName, person.LastName))
		email := strings.ToLower(faker.Email())

		password := faker.Password(true, true, true, false, false, 12)
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}

		createdAt := now.Add(-time.Duration(faker.Number(0, 365)) * 24 * time.Hour)
		updatedAt := createdAt
		if delta := faker.Number(0, int(time.Since(createdAt).Hours())); delta > 0 {
			updatedAt = createdAt.Add(time.Duration(delta) * time.Hour)
			if updatedAt.After(now) {
				updatedAt = now
			}
		}

		var lastLogin *time.Time
		if faker.Bool() {
			candidate := createdAt.Add(time.Duration(faker.Number(1, 240)) * time.Hour)
			if candidate.After(now) {
				candidate = now
			}
			lastLogin = &candidate
		}

		var avatar *string
		if faker.Bool() {
			seed := url.QueryEscape(strings.ReplaceAll(name, " ", ""))
			value := fmt.Sprintf("https://api.dicebear.com/7.x/initials/svg?seed=%s", seed)
			avatar = &value
		}

		users[i] = User{
			UUID:          uuid.NewString(),
			Name:          name,
			Email:         email,
			PasswordHash:  string(hash),
			AvatarURL:     avatar,
			Status:        statusPool[faker.Number(0, len(statusPool)-1)],
			LastLoginAt:   lastLogin,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			PlainPassword: password,
		}
	}

	return users, nil
}
