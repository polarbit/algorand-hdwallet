package bip32path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFails(t *testing.T) {

	var invalidPaths = []string{
		"s",
		"/0",
		"M/0",
		"m/0''",
		"m/0/",
		"m/0//1",
		"m/0/1/2/3/4/5/6",
		"m/3000000000",
		"m/01/1",
		"m/0/s",
	}

	for _, path := range invalidPaths {
		t.Run(path, func(t *testing.T) {
			_, err := Parse(path)
			assert.Error(t, err)
		})
	}
}

func TestParseSucceeds(t *testing.T) {

	validPaths := []struct {
		rawPath   string
		bip32Path *Bip32Path
	}{
		{
			rawPath: "m",
			bip32Path: &Bip32Path{
				Segments: nil,
			},
		},
		{
			rawPath: "m/0'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/1",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 1, IsHardened: false},
				},
			},
		},
		{
			rawPath: "m/0'/1'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: true},
					{Value: 1, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/0'/1'/2'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: true},
					{Value: 1, IsHardened: true},
					{Value: 2, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/0'/1'/2'/2'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: true},
					{Value: 1, IsHardened: true},
					{Value: 2, IsHardened: true},
					{Value: 2, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/0'/1'/2'/2'/1000000000'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: true},
					{Value: 1, IsHardened: true},
					{Value: 2, IsHardened: true},
					{Value: 2, IsHardened: true},
					{Value: 1000000000, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/0",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: false},
				},
			},
		},
		{
			rawPath: "m/0/2147483647'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: false},
					{Value: 2147483647, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/0/2147483647'/1",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: false},
					{Value: 2147483647, IsHardened: true},
					{Value: 1, IsHardened: false},
				},
			},
		},
		{
			rawPath: "m/0/2147483647'/1/2147483646'",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: false},
					{Value: 2147483647, IsHardened: true},
					{Value: 1, IsHardened: false},
					{Value: 2147483646, IsHardened: true},
				},
			},
		},
		{
			rawPath: "m/0/2147483647'/1/2147483646'/2",
			bip32Path: &Bip32Path{
				Segments: []*Bip32PathSegment{
					{Value: 0, IsHardened: false},
					{Value: 2147483647, IsHardened: true},
					{Value: 1, IsHardened: false},
					{Value: 2147483646, IsHardened: true},
					{Value: 2, IsHardened: false},
				},
			},
		},
	}

	for _, testPath := range validPaths {
		t.Run(testPath.rawPath, func(t *testing.T) {
			parsed, err := Parse(testPath.rawPath)
			assert.NoError(t, err)
			assert.Equal(t, testPath.rawPath, parsed.RawPath)

			for i, xpart := range testPath.bip32Path.Segments {
				tpart := parsed.Segments[i]

				assert.Equal(t, xpart.Value, tpart.Value)
				assert.Equal(t, xpart.IsHardened, tpart.IsHardened)
			}

		})
	}
}
