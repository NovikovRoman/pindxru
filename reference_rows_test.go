package pindxru

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_GetLastModified(t *testing.T) {
	d, err := testReferenceRows.GetLastModified()
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.False(t, d.IsZero())
}
