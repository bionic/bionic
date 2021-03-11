package imports

import (
	"errors"
	"fmt"
	"github.com/bionic-dev/bionic/database"
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	manager, err := NewManager(db, []provider.Provider{p})
	require.NoError(t, err)

	assert.Equal(t, map[string]provider.Provider{"mock": p}, manager.providers)
}

func TestManager_Migrate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := provider.NewMockProvider(ctrl)

	p.EXPECT().
		Name().
		Return("mock")
	p.EXPECT().
		Migrate().
		Return(nil)

	manager, err := NewManager(db, []provider.Provider{p})
	require.NoError(t, err)

	err = manager.Migrate()
	require.NoError(t, err)
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

func TestDefaultProviders(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	providers := DefaultProviders(db)

	var tables []string

	initialTables, err := database.GetTables(db)
	require.Nil(t, err)
	tables = initialTables

	t.Run("model tables have prefixes", func(t *testing.T) {
		for _, p := range providers {
			err := p.Migrate()
			require.NoError(t, err)

			currentTables, err := database.GetTables(db)
			require.NoError(t, err)

			addedTables := sliceDiff(currentTables, tables)
			for _, table := range addedTables {
				assert.True(t, strings.HasPrefix(table, p.TablePrefix()),
					fmt.Sprintf(`"%s" table does not have "%s" prefix`, table, p.TablePrefix()))
			}

			tables = currentTables
		}
	})
}

func sliceDiff(slice1, slice2 []string) []string {
	sliceMapping := make(map[string]bool, len(slice2))
	for _, x := range slice2 {
		sliceMapping[x] = true
	}

	var diff []string
	for _, value := range slice1 {
		if !sliceMapping[value] {
			diff = append(diff, value)
		}
	}

	return diff
}
