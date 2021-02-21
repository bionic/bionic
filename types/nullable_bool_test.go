package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNullableBool_UnmarshalCSV(t *testing.T) {
	nb := &NullableBool{}

	t.Run("true value", func(t *testing.T) {
		err := nb.UnmarshalCSV("1")
		require.NoError(t, err)
		assert.Equal(t, true, nb.Valid)
		assert.Equal(t, true, nb.Bool)
	})

	t.Run("false value", func(t *testing.T) {
		err := nb.UnmarshalCSV("0")
		require.NoError(t, err)
		assert.Equal(t, true, nb.Valid)
		assert.Equal(t, false, nb.Bool)
	})

	t.Run("invalid value", func(t *testing.T) {
		err := nb.UnmarshalCSV("-1")
		require.NoError(t, err)
		assert.Equal(t, false, nb.Valid)
	})

	t.Run("empty value", func(t *testing.T) {
		err := nb.UnmarshalCSV("")
		require.NoError(t, err)
		assert.Equal(t, false, nb.Valid)
	})
}
