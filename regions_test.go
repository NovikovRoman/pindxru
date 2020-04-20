package pindxru

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegions_GetCode(t *testing.T) {
	code, autonomy := Regions.GetCode("ИРКУТСКАЯ область")
	require.Equal(t, code, 38)
	require.False(t, autonomy)

	code, autonomy = Regions.GetCode("ХАНТЫ-МАНСИЙСКИЙ-югра АВТОНОМНЫЙ округ")
	require.Equal(t, code, 86)
	require.True(t, autonomy)

	code, autonomy = Regions.GetCode("unknown")
	require.Equal(t, code, 0)
	require.False(t, autonomy)
}

func TestRegions_GetName(t *testing.T) {
	name := Regions.GetName(38)
	require.Equal(t, name, "ИРКУТСКАЯ ОБЛАСТЬ")

	name = Regions.GetName(86)
	require.Equal(t, name, "ХАНТЫ-МАНСИЙСКИЙ-ЮГРА АВТОНОМНЫЙ ОКРУГ")

	name = Regions.GetName(0)
	require.Equal(t, name, "")

	name = Regions.GetName(-100)
	require.Equal(t, name, "")

	name = Regions.GetName(100)
	require.Equal(t, name, "")
}

func TestRegions_FindCode(t *testing.T) {
	code := Regions.FindCode("Тюменская область", "ХАНТЫ-МАНСИЙСКИЙ-ЮГРА АВТОНОМНЫЙ ОКРУГ")
	require.Equal(t, code, 72)

	code = Regions.FindCode("", "ХАНТЫ-МАНСИЙСКИЙ-ЮГРА АВТОНОМНЫЙ ОКРУГ")
	require.Equal(t, code, 86)

	code = Regions.FindCode("", "")
	require.Equal(t, code, 0)
}

func TestRegions_FindRegionCodeByIndex(t *testing.T) {
	code, err := FindRegionCodeByIndex("Тюменская область")
	require.NotNil(t, err)
	require.Equal(t, code, 0)

	code, err = FindRegionCodeByIndex("165")
	require.Nil(t, err)
	require.Equal(t, code, 29)
}
