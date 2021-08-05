package pindxru

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_getListUpdates(t *testing.T) {
	var (
		err error
		b   []byte
		u   []Package
	)
	b, err = ioutil.ReadFile(filepath.Join(testdata, "page.htm"))
	require.Nil(t, err)
	u, err = getListUpdates(b)
	require.Nil(t, err)
	require.Len(t, u, 14)
}
