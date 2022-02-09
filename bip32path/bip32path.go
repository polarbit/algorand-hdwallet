package bip32path

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	HARDENED_OFFSET = 0x80000000
)

type Bip32Path struct {
	// Purpose      uint32
	// CoinType     uint32
	// Account      uint32
	// Change       uint32
	// AddressIndex uint32
	// Depth        uint8
	RawPath  string
	Segments []*Bip32PathSegment
}

type Bip32PathSegment struct {
	Value      uint32
	IsHardened bool
}

var r *regexp.Regexp

func init() {
	r, _ = regexp.Compile("^m(\\/[0-9]+'{0,1}){0,5}$")
}

func Parse(s string) (*Bip32Path, error) {
	if false == r.MatchString(s) {
		return nil, errors.New("Path is invalid (1)")
	}

	segments := strings.Split(s[1:], "/")

	path := &Bip32Path{
		RawPath:  s,
		Segments: make([]*Bip32PathSegment, 0),
	}

	for i, s := range segments {
		if i == 0 {
			continue
		}

		newseg := Bip32PathSegment{}
		path.Segments = append(path.Segments, &newseg)

		if s[len(s)-1] == byte('\'') {
			val, err := strconv.ParseUint(s[:len(s)-1], 10, 32)
			if err != nil {
				return nil, errors.New("Path is invalid (2)")
			}
			newseg.Value = uint32(val)
			newseg.IsHardened = true
		} else {
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil, errors.New("Path is invalid (3)")
			}
			newseg.Value = uint32(val)
		}

		if newseg.Value > 0 && s[0] == byte('0') {
			return nil, errors.New("Path is invalid (4)")
		}

		if newseg.Value >= HARDENED_OFFSET {
			return nil, errors.New("Path is invalid (5)")
		}
	}

	return path, nil
}
