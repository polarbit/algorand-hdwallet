package utils

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/ripemd160"
)

func HashSha256(data []byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func HashRipeMD160(data []byte) ([]byte, error) {
	hasher := ripemd160.New()
	_, err := io.WriteString(hasher, string(data))
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func Hash160(data []byte) ([]byte, error) {
	hash1, err := HashSha256(data)
	if err != nil {
		return nil, err
	}

	hash2, err := HashRipeMD160(hash1)
	if err != nil {
		return nil, err
	}

	return hash2, nil
}
