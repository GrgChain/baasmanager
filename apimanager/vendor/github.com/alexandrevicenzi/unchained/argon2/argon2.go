package argon2

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Errors returned by Argon2Hasher.
var (
	ErrHashComponentUnreadable = errors.New("unchained/argon2: unreadable component in hashed password")
	ErrHashComponentMismatch   = errors.New("unchained/argon2: hashed password components mismatch")
	ErrAlgorithmMismatch       = errors.New("unchained/argon2: algorithm mismatch")
	ErrIncompatibleVersion     = errors.New("unchained/argon2: incompatible version")
)

// Argon2Hasher implements Argon2i password hasher.
type Argon2Hasher struct {
	// Algorithm identifier.
	Algorithm string
	// Defines the amount of computation time, given in number of iterations.
	Time uint32
	// Defines the memory usage (KiB).
	Memory uint32
	// Defines the number of parallel threads.
	Threads uint8
	// Defines the length of the hash in bytes.
	Length uint32
}

// Encode turns a plain-text password into a hash.
func (h *Argon2Hasher) Encode(password string, salt string) (string, error) {
	bSalt := []byte(salt)
	hash := argon2.Key([]byte(password), bSalt, h.Time, h.Memory, h.Threads, h.Length)

	b64Salt := base64.RawStdEncoding.EncodeToString(bSalt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	s := fmt.Sprintf("%s$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		h.Algorithm,
		"argon2i",
		argon2.Version,
		h.Memory,
		h.Time,
		h.Threads,
		b64Salt,
		b64Hash,
	)

	return s, nil
}

// Verify if a plain-text password matches the encoded digest.
func (h *Argon2Hasher) Verify(password string, encoded string) (bool, error) {
	s := strings.Split(encoded, "$")

	if len(s) != 6 {
		return false, ErrHashComponentMismatch
	}

	algorithm, method, version, params, salt, hash := s[0], s[1], s[2], s[3], s[4], s[5]

	if algorithm != h.Algorithm || method != "argon2i" {
		return false, ErrAlgorithmMismatch
	}

	var v int
	var err error

	_, err = fmt.Sscanf(version, "v=%d", &v)

	if err != nil {
		return false, ErrHashComponentUnreadable
	}

	if v != argon2.Version {
		return false, ErrIncompatibleVersion
	}

	var time, memory uint32
	var threads uint8

	_, err = fmt.Sscanf(params, "m=%d,t=%d,p=%d", &memory, &time, &threads)

	if err != nil {
		return false, ErrHashComponentUnreadable
	}

	bSalt, err := base64.RawStdEncoding.DecodeString(salt)

	if err != nil {
		return false, ErrHashComponentUnreadable
	}

	bHash, err := base64.RawStdEncoding.DecodeString(hash)

	if err != nil {
		return false, ErrHashComponentUnreadable
	}

	newHash := argon2.Key([]byte(password), bSalt, time, memory, threads, uint32(len(bHash)))

	return subtle.ConstantTimeCompare(bHash, newHash) == 1, nil
}

// NewArgon2Hasher secures password hashing using the argon2 algorithm.
func NewArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{
		Algorithm: "argon2",
		Time:      2,
		Memory:    512,
		Threads:   2,
		Length:    16,
	}
}
