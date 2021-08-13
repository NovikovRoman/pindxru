package pindxru

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_Indexes(t *testing.T) {
	u, d, err := cTest.Indexes(testReferenceRows, nil)
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.False(t, d.IsZero())
	require.True(t, len(u) > 10000)

	// (!) обновлений нет
	lastMod := time.Now().Add(time.Hour * 1000)
	u, d, err = cTest.Indexes(testReferenceRows, &lastMod)
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.True(t, d.IsZero())
	require.Len(t, u, 0)

	lastMod = time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	u, d, err = cTest.Indexes(testReferenceRows, &lastMod)
	require.Nil(t, err)
	require.IsType(t, d, time.Time{})
	require.False(t, d.IsZero())
	require.True(t, len(u) > 10000)
}

func Test_IndexesZip(t *testing.T) {
	filename := filepath.Join(testdata, "indexes-"+testZipFile)
	lastMod, ok, err := cTest.IndexesZip(testReferenceRows, filename, os.ModePerm, nil)
	require.Nil(t, err)
	require.True(t, ok)
	require.False(t, lastMod.IsZero())

	testCheckFile(t, filename)

	filename = filepath.Join(testdata, "indexesWithDate-"+testZipFile)
	// (!) обновлений нет
	d := time.Now().Add(time.Hour * 1000)
	lastMod, ok, err = cTest.IndexesZip(testReferenceRows, filename, os.ModePerm, &d)
	require.Nil(t, err)
	require.False(t, ok)
	require.True(t, lastMod.IsZero())

	d = time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	lastMod, ok, err = cTest.IndexesZip(testReferenceRows, filename, os.ModePerm, &d)
	require.Nil(t, err)
	require.True(t, ok)
	require.False(t, lastMod.IsZero())

	testCheckFile(t, filename)
}

func Test_IndexesDbf(t *testing.T) {
	filename := filepath.Join(testdata, "indexes-"+testDbfFile)
	lastMod, ok, err := cTest.IndexesDbf(testReferenceRows, filename, os.ModePerm, nil)
	require.Nil(t, err)
	require.True(t, ok)
	require.False(t, lastMod.IsZero())

	testCheckFile(t, filename)

	filename = filepath.Join(testdata, "indexesWithDate-"+testDbfFile)
	// (!) обновлений нет
	d := time.Now().Add(time.Hour * 1000)
	lastMod, ok, err = cTest.IndexesDbf(testReferenceRows, filename, os.ModePerm, &d)
	require.Nil(t, err)
	require.False(t, ok)
	require.True(t, lastMod.IsZero())

	d = time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	lastMod, ok, err = cTest.IndexesDbf(testReferenceRows, filename, os.ModePerm, &d)
	require.Nil(t, err)
	require.True(t, ok)
	require.False(t, lastMod.IsZero())

	testCheckFile(t, filename)
}

func Test_Packages(t *testing.T) {
	// (!) обновлений нет
	d := time.Now().Add(time.Hour * 1000)
	pack, err := testReferenceRows.GetUpdatePackages(&d)
	require.Nil(t, err)
	require.Len(t, pack, 0)

	lastMod, err := cTest.GetPackageIndexes(&testPackages[0])
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())
	require.True(t, len(testPackages[0].Indexes) > 0)
}

func Test_PackageZip(t *testing.T) {
	filename := filepath.Join(testdata, "package-"+testZipFile)
	lastMod, err := cTest.PackageZip(testPackages[0], filename, os.ModePerm)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())

	testCheckFile(t, filename)
}

func Test_PackageDbf(t *testing.T) {
	filename := filepath.Join(testdata, "package-"+testDbfFile)
	lastMod, err := cTest.PackageDbf(testPackages[0], filename, os.ModePerm)
	require.Nil(t, err)
	require.False(t, lastMod.IsZero())

	testCheckFile(t, filename)
}

func testCheckFile(t *testing.T, filename string) {
	f, err := os.Open(filename)
	require.Nil(t, err)
	fi, err := f.Stat()
	require.Nil(t, err)
	require.True(t, fi.Size() > 0)

	require.Nil(t, os.Remove(filename))
}
