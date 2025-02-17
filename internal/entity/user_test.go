package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "test@test.com", "123321")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "test@test.com", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("John Doe", "test@test.com", "123321")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123321"))
	assert.False(t, user.ValidatePassword("123456"))
	assert.NotEqual(t, "123321", user.Password)
}
