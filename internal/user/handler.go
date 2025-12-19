package user

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kouji-app-backend2/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler wires user-related HTTP endpoints.
type Handler struct {
	db *gorm.DB
}

// NewHandler returns a Handler instance.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Register attaches the HTTP routes to the router.
func (h *Handler) Register(router *echo.Echo) {
	router.GET("/users", h.ListUsers)
	router.GET("/users/:id", h.GetUser)
}

// ListUsers returns paginated users ordered by newest first.
func (h *Handler) ListUsers(c echo.Context) error {
	limit := parsePositiveInt(c.QueryParam("limit"), 50)
	if limit > 100 {
		limit = 100
	}

	page := parsePositiveInt(c.QueryParam("page"), 1)
	offset := (page - 1) * limit

	var users []models.User
	if err := h.db.Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("failed to query users"))
	}

	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = newUserResponse(u)
	}

	return c.JSON(http.StatusOK, UsersResponse{
		Data: resp,
		Meta: PaginationMeta{Page: page, Limit: limit, Count: len(resp)},
	})
}

// GetUser returns a specific user identified by numeric ID or UUID.
func (h *Handler) GetUser(c echo.Context) error {
	identifier := strings.TrimSpace(c.Param("id"))
	if identifier == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("user identifier is required"))
	}

	var user models.User
	var err error

	if id, parseErr := strconv.ParseUint(identifier, 10, 64); parseErr == nil {
		err = h.db.First(&user, id).Error
	} else {
		err = h.db.Where("uuid = ?", identifier).First(&user).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, errorResponse("user not found"))
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("failed to query user"))
	}

	return c.JSON(http.StatusOK, newUserResponse(user))
}

// UserResponse describes the public fields returned to API callers.
type UserResponse struct {
	ID          uint64            `json:"id"`
	UUID        string            `json:"uuid"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	AvatarURL   *string           `json:"avatar_url,omitempty"`
	Status      models.UserStatus `json:"status"`
	LastLoginAt *time.Time        `json:"last_login_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// UsersResponse wraps a list of users together with simple paging info.
type UsersResponse struct {
	Data []UserResponse `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta conveys paging details for list responses.
type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Count int `json:"count"`
}

func newUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		UUID:        user.UUID,
		Name:        user.Name,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func errorResponse(message string) map[string]string {
	return map[string]string{"error": message}
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
