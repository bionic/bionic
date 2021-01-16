package providers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/shekhirin/bionic-cli/database"
	"github.com/shekhirin/bionic-cli/providers/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
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

	p.EXPECT().
		Models().
		Return([]schema.Tabler{&database.MockModel{}})

	manager, err := NewManager(db, []provider.Provider{p})
	require.NoError(t, err)

	assert.Equal(t, map[string]provider.Provider{"mock": p}, manager.providers)

	assert.True(t, db.Migrator().HasTable(&database.MockModel{}))
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

		p.EXPECT().
			Models().
			Return([]schema.Tabler{&database.MockModel{}})

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

func TestDefaultProviders(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	providers := DefaultProviders(db)

	t.Run("model tables have prefixes", func(t *testing.T) {
		for _, p := range providers {
			for _, model := range p.Models() {
				assert.True(t, strings.HasPrefix(model.TableName(), p.TablePrefix()))
			}
		}
	})
}
