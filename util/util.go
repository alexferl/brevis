package util

import (
	"crypto/sha256"
	"net/url"
)

// GetSha256Hash hashes a string and returns a sha256 byte slice
func GetSha256Hash(text string) [32]byte {
	return sha256.Sum256([]byte(text))
}

// IsValidUri checks whether a URI is valid or not
func IsValidUri(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	return true
}
