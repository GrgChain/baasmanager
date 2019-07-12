package bcrypt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Errors returned by BCryptHasher.
var (
	ErrHashComponentMismatch = errors.New("unchained/bcrypt: hashed password components mismatch")
	ErrAlgorithmMismatch     = errors.New("unchained/bcrypt: algorithm mismatch")
)

// BCryptHasher implements Bcrypt password hasher.
type BCryptHasher struct {
	// Algorithm identifier.
	Algorithm string
	// Defines the hash function used to avoid bcrypt's 72 bytes password truncation.
	Digest func() hash.Hash
	// Defines the number of rounds used to encode the password.
	Cost int
}

// Encode turns a plain-text password into a hash.
//
// Parameter salt is currently ignored.
func (h *BCryptHasher) Encode(password string, salt string) (string, error) {
	if h.Digest != nil {
		d := h.Digest()
		d.Write([]byte(password))
		password = hex.EncodeToString(d.Sum(nil))
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s$%s", h.Algorithm, string(bytes)), nil
}

// Verify if a plain-text password matches the encoded digest.
func (h *BCryptHasher) Verify(password string, encoded string) (bool, error) {
	s := strings.SplitN(encoded, "$", 2)

	if len(s) != 2 {
		return false, ErrHashComponentMismatch
	}

	algorithm, hash := s[0], s[1]

	if algorithm != h.Algorithm {
		return false, ErrAlgorithmMismatch
	}

	if h.Digest != nil {
		d := h.Digest()
		d.Write([]byte(password))
		password = hex.EncodeToString(d.Sum(nil))
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, nil
}

// NewBCryptHasher secures password hashing using the bcrypt algorithm.
//
// This hasher does not first hash the password which means it is subject to
// bcrypt's 72 bytes password truncation.
func NewBCryptHasher() *BCryptHasher {
	return &BCryptHasher{
		Algorithm: "bcrypt",
		Digest:    nil,
		Cost:      12,
	}
}

// NewBCryptSHA256Hasher secures password hashing using the bcrypt algorithm.
//
// This hasher first hash the password with SHA-256.
func NewBCryptSHA256Hasher() *BCryptHasher {
	return &BCryptHasher{
		Algorithm: "bcrypt_sha256",
		Digest:    sha256.New,
		Cost:      12,
	}
}
