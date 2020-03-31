package pindxru

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func Test_GetLastModified(t *testing.T) {
	d, err := cTest.GetLastModified()
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.False(t, d.IsZero())
}

func Test_All(t *testing.T) {
	u, d, err := cTest.All()
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.False(t, d.IsZero())
	require.True(t, len(u) > 10000)
}

func Test_AllUpdates(t *testing.T) {
	// (!) обновлений нет
	u, d, err := cTest.AllUpdated(time.Now().Add(time.Hour * 1000))
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.True(t, d.IsZero())
	require.Len(t, u, 0)

	u, d, err = cTest.AllUpdated(time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC))
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.False(t, d.IsZero())
	require.True(t, len(u) > 10000)
}

func Test_Updates(t *testing.T) {
	// (!) обновлений нет
	d := time.Now().Add(time.Hour * 1000)
	u, err := cTest.Updates(&d)
	require.Nil(t, err)
	require.Len(t, u, 0)

	d = time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	u, err = cTest.Updates(&d)
	require.Nil(t, err)
	require.True(t, len(u) > 0)

	for _, i := range u {
		require.False(t, i.Date.IsZero())
		require.True(t, i.NumberRecords > 0)
	}
}

func Test_GetNPIndxes(t *testing.T) {
	d := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	u, err := cTest.Updates(&d)
	require.Nil(t, err)

	indexes, lastMod, err := cTest.GetNPIndxes(u[0].Url)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())
	require.True(t, len(indexes) > 0)
}

func Test_ZipAll(t *testing.T) {
	filename := "zipall-" + testZipFile
	lastMod, err := cTest.ZipAll(filename, os.ModePerm)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())

	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	err = os.Remove(filename)
	require.Nil(t, err)
}

func Test_ZipAllUpdates(t *testing.T) {
	filename := "zipallupdates-" + testZipFile
	// (!) обновлений нет
	d := time.Now().Add(time.Hour * 1000)
	lastMod, ok, err := cTest.ZipAllUpdated(filename, os.ModePerm, d)
	require.Nil(t, err)
	require.False(t, ok)
	require.True(t, lastMod.IsZero())

	d = time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	lastMod, ok, err = cTest.ZipAllUpdated(filename, os.ModePerm, d)
	require.Nil(t, err)
	require.True(t, ok)
	require.False(t, lastMod.IsZero())

	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	err = os.Remove(filename)
	require.Nil(t, err)
}

func Test_DbfAll(t *testing.T) {
	filename := "dbfall-" + testDbfFile
	lastMod, err := cTest.DbfAll(filename, os.ModePerm)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())

	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	err = os.Remove(filename)
	require.Nil(t, err)
}

func Test_DbfAllUpdates(t *testing.T) {
	filename := "dbfallupdates-" + testDbfFile
	// (!) обновлений нет
	d := time.Now().Add(time.Hour * 1000)
	lastMod, ok, err := cTest.DbfAllUpdated(filename, os.ModePerm, d)
	require.Nil(t, err)
	require.False(t, ok)
	require.True(t, lastMod.IsZero())

	d = time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	lastMod, ok, err = cTest.DbfAllUpdated(filename, os.ModePerm, d)
	require.Nil(t, err)
	require.True(t, ok)
	require.False(t, lastMod.IsZero())

	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	err = os.Remove(filename)
	require.Nil(t, err)
}

func Test_ZipNPIndx(t *testing.T) {
	d := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	u, err := cTest.Updates(&d)
	require.Nil(t, err)

	filename := "zipNPIndx-" + testZipFile
	lastMod, err := cTest.ZipNPIndx(u[0].Url, filename, os.ModePerm)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())

	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	err = os.Remove(filename)
	require.Nil(t, err)
}

func Test_DbfNPIndx(t *testing.T) {
	d := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	u, err := cTest.Updates(&d)
	require.Nil(t, err)

	filename := "dbfNPIndx-" + testZipFile
	lastMod, err := cTest.DbfNPIndx(u[0].Url, filename, os.ModePerm)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())

	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	err = os.Remove(filename)
	require.Nil(t, err)
}
