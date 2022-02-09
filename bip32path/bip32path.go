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
	RawPath string
	Parts   []*Bip32PathPart
}

type Bip32PathPart struct {
	Value      uint32
	IsHardened bool
}

var r *regexp.Regexp

func init() {
	r, _ = regexp.Compile("^m(\\/[0-9]+'{0,1}){0,5}$")
}

func IsValid(s string) bool {
	return r.MatchString(s)
}

func Parse(s string) (*Bip32Path, error) {
	if false == r.MatchString(s) {
		return nil, errors.New("Path is invalid (1)")
	}

	parts := strings.Split(s[1:], "/")

	path := &Bip32Path{
		RawPath: s,
		Parts:   make([]*Bip32PathPart, 0),
	}

	for i, p := range parts {
		if i == 0 {
			continue
		}

		newp := Bip32PathPart{}
		path.Parts = append(path.Parts, &newp)

		if p[len(p)-1] == byte('\'') {
			val, err := strconv.ParseUint(p[:len(p)-1], 10, 32)
			if err != nil {
				return nil, errors.New("Path is invalid (2)")
			}
			newp.Value = uint32(val)
			newp.IsHardened = true
		} else {
			val, err := strconv.ParseUint(p, 10, 32)
			if err != nil {
				return nil, errors.New("Path is invalid (3)")
			}
			newp.Value = uint32(val)
		}

		if newp.Value > 0 && p[0] == byte('0') {
			return nil, errors.New("Path is invalid (4)")
		}

		if newp.Value >= HARDENED_OFFSET {
			return nil, errors.New("Path is invalid (5)")
		}
	}

	return path, nil
}
