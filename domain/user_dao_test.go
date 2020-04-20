package domain

/**
These are gonna be simple test, because so far, we were using a mock database. But a really
important thing to notice is the usage of 'testify' library. If an assertion failed, it will continue
executing the others.
 */
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserNoDataFound(t *testing.T) {

	user, err := GetUser(0)
	assert.Nil(t, user, "user should be null due non existent data")
	assert.NotNil(t, err, "Error message should not be null")
}

func TestGetUserFound(t *testing.T) {

	user, err := GetUser(1)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@domain.com", user.Email)
	assert.Equal(t, "TestFirstName", user.FirstName)
	assert.Equal(t, "TestLastName", user.LastName)
}
