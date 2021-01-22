package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNullableBool_UnmarshalCSV(t *testing.T) {
	nb := &NullableBool{}

	err := nb.UnmarshalCSV("1")
	require.NoError(t, err)
	assert.Equal(t, true, nb.Valid)
	assert.Equal(t, true, nb.Bool)

	err = nb.UnmarshalCSV("0")
	require.NoError(t, err)
	assert.Equal(t, true, nb.Valid)
	assert.Equal(t, false, nb.Bool)

	err = nb.UnmarshalCSV("-1")
	require.NoError(t, err)
	assert.Equal(t, false, nb.Valid)

	err = nb.UnmarshalCSV("")
	require.NoError(t, err)
	assert.Equal(t, false, nb.Valid)
}
