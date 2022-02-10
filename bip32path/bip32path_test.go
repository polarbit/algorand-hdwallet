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
			"m",
			&Bip32Path{
				Segments: nil,
			},
		},
		{
			"m/0'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: true},
				},
			},
		},
		{
			"m/1",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 1, IsHardened: false},
				},
			},
		},
		{
			"m/0'/1'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: true},
					{ValueSeen: 1, IsHardened: true},
				},
			},
		},
		{
			"m/0'/1'/2'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: true},
					{ValueSeen: 1, IsHardened: true},
					{ValueSeen: 2, IsHardened: true},
				},
			},
		},
		{
			"m/0'/1'/2'/2'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: true},
					{ValueSeen: 1, IsHardened: true},
					{ValueSeen: 2, IsHardened: true},
					{ValueSeen: 2, IsHardened: true},
				},
			},
		},
		{
			"m/0'/1'/2'/2'/1000000000'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: true},
					{ValueSeen: 1, IsHardened: true},
					{ValueSeen: 2, IsHardened: true},
					{ValueSeen: 2, IsHardened: true},
					{ValueSeen: 1000000000, IsHardened: true},
				},
			},
		},
		{
			"m/0",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: false},
				},
			},
		},
		{
			"m/0/2147483647'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: false},
					{ValueSeen: 2147483647, IsHardened: true},
				},
			},
		},
		{
			"m/0/2147483647'/1",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: false},
					{ValueSeen: 2147483647, IsHardened: true},
					{ValueSeen: 1, IsHardened: false},
				},
			},
		},
		{
			"m/0/2147483647'/1/2147483646'",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: false},
					{ValueSeen: 2147483647, IsHardened: true},
					{ValueSeen: 1, IsHardened: false},
					{ValueSeen: 2147483646, IsHardened: true},
				},
			},
		},
		{
			"m/0/2147483647'/1/2147483646'/2",
			&Bip32Path{
				Segments: []*Bip32PathSegment{
					{ValueSeen: 0, IsHardened: false},
					{ValueSeen: 2147483647, IsHardened: true},
					{ValueSeen: 1, IsHardened: false},
					{ValueSeen: 2147483646, IsHardened: true},
					{ValueSeen: 2, IsHardened: false},
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

				assert.Equal(t, xpart.ValueSeen, tpart.ValueSeen)
				assert.Equal(t, xpart.IsHardened, tpart.IsHardened)
			}

		})
	}
}
