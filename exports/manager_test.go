package exports

import (
	"errors"
	"github.com/bionic-dev/bionic/exports/provider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestNewManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := provider.NewMockProvider(ctrl)

	p.EXPECT().
		Name().
		Return("mock")

	manager, err := NewManager(db, []provider.Provider{p})
	require.NoError(t, err)

	assert.Equal(t, map[string]provider.Provider{"mock": p}, manager.providers)
}

func TestManager_GetByName(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		require.NoError(t, err)

		p := provider.NewMockProvider(ctrl)

		p.EXPECT().
			Name().
			Return("mock")

		manager, err := NewManager(db, []provider.Provider{p})
		require.NoError(t, err)

		pByName, err := manager.GetByName("mock")
		require.NoError(t, err)

		assert.Equal(t, p, pByName)
	})

	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		require.NoError(t, err)

		manager, err := NewManager(db, []provider.Provider{})
		require.NoError(t, err)

		pByName, err := manager.GetByName("mock")
		require.Nil(t, pByName)

		assert.True(t, errors.Is(err, ErrProviderNotFound))
	})
}
