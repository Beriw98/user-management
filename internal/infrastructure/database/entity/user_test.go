package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Beriw98/user-management/internal/infrastructure/database/entity"
)

func TestUser_Fields(t *testing.T) {
	t.Run("Fields", func(t *testing.T) {
		u := entity.User{}
		got := u.Fields()

		assert.Len(t, got, 5)
		assert.Equal(t, "id", got[0].Descriptor().Name)
		assert.Equal(t, "name", got[1].Descriptor().Name)
		assert.Equal(t, "surname", got[2].Descriptor().Name)
		assert.Equal(t, "email", got[3].Descriptor().Name)
		assert.Equal(t, "password", got[4].Descriptor().Name)
	})
}
