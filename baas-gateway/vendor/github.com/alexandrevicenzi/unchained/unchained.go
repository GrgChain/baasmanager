package unchained

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/alexandrevicenzi/unchained/argon2"
	"github.com/alexandrevicenzi/unchained/bcrypt"
	"github.com/alexandrevicenzi/unchained/pbkdf2"
)

// Django hasher identifiers.
const (
	Argon2Hasher       = "argon2"
	BCryptHasher       = "bcrypt"
	BCryptSHA256Hasher = "bcrypt_sha256"
	CryptHasher        = "crypt"
	MD5Hasher          = "md5"
	PBKDF2SHA1Hasher   = "pbkdf2_sha1"
	PBKDF2SHA256Hasher = "pbkdf2_sha256"
	SHA1Hasher         = "sha1"
	UnsaltedMD5Hasher  = "unsalted_md5"
	UnsaltedSHA1Hasher = "unsalted_sha1"
)

const (
	// The prefix used in unusable passwords.
	UnusablePasswordPrefix = "!"
	// The length of unusable passwords after the prefix.
	UnusablePasswordSuffixLength = 40
	// The default hasher used in Django.
	DefaultHasher = PBKDF2SHA256Hasher
	// The default salt size used in Django.
	DefaultSaltSize = 12

	allowedChars     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	allowedCharsSize = len(allowedChars)
)

var (
	// ErrInvalidHasher is returned if the hasher is invalid or unknown.
	ErrInvalidHasher = errors.New("unchained: invalid hasher")
	// ErrHasherNotImplemented is returned if the hasher is not implemented.
	ErrHasherNotImplemented = errors.New("unchained: hasher not implemented")
)

// GetRandomString returns a securely generated random string.
func GetRandomString(length int) string {
	b := make([]byte, length)

	for i := range b {
		c := rand.Intn(allowedCharsSize)
		b[i] = allowedChars[c]
	}

	return string(b)
}

// IsValidHasher returns true if the hasher
// is supported by Django, or false otherwise.
func IsValidHasher(hasher string) bool {
	switch hasher {
	case
		Argon2Hasher,
		BCryptHasher,
		BCryptSHA256Hasher,
		CryptHasher,
		MD5Hasher,
		PBKDF2SHA1Hasher,
		PBKDF2SHA256Hasher,
		SHA1Hasher,
		UnsaltedMD5Hasher,
		UnsaltedSHA1Hasher:
		return true
	}

	return false
}

// IsWeakHasher returns true if the hasher
// is not recommend by Django, or false otherwise.
func IsWeakHasher(hasher string) bool {
	switch hasher {
	case
		CryptHasher,
		MD5Hasher,
		SHA1Hasher,
		UnsaltedMD5Hasher,
		UnsaltedSHA1Hasher:
		return true
	}

	return false
}

// IsHasherImplemented returns true if the hasher
// is implemented in this library, or false otherwise.
func IsHasherImplemented(hasher string) bool {
	switch hasher {
	case
		Argon2Hasher,
		BCryptHasher,
		BCryptSHA256Hasher,
		PBKDF2SHA1Hasher,
		PBKDF2SHA256Hasher:
		return true
	}

	return false
}

// IdentifyHasher returns the hasher used in the encoded password.
func IdentifyHasher(encoded string) string {
	size := len(encoded)

	if size == 32 && !strings.Contains(encoded, "$") {
		return UnsaltedMD5Hasher
	}

	if size == 37 && strings.HasPrefix(encoded, "md5$$") {
		return UnsaltedMD5Hasher
	}

	if size == 46 && strings.HasPrefix(encoded, "sha1$$") {
		return UnsaltedSHA1Hasher
	}

	return strings.SplitN(encoded, "$", 2)[0]
}

// IsPasswordUsable returns true if encoded password
// is usable, or false otherwise.
func IsPasswordUsable(encoded string) bool {
	return encoded != "" && !strings.HasPrefix(encoded, UnusablePasswordPrefix)
}

// CheckPassword validate if the raw password matches the encoded digest.
//
// This is a shortcut that discovers the hasher used in the encoded digest
// to perform the correct validation.
func CheckPassword(password, encoded string) (bool, error) {
	if !IsPasswordUsable(encoded) {
		return false, nil
	}

	hasher := IdentifyHasher(encoded)

	if !IsValidHasher(hasher) {
		return false, ErrInvalidHasher
	}

	if !IsHasherImplemented(hasher) {
		return false, ErrHasherNotImplemented
	}

	switch hasher {
	case Argon2Hasher:
		return argon2.NewArgon2Hasher().Verify(password, encoded)
	case BCryptHasher:
		return bcrypt.NewBCryptHasher().Verify(password, encoded)
	case BCryptSHA256Hasher:
		return bcrypt.NewBCryptSHA256Hasher().Verify(password, encoded)
	case PBKDF2SHA1Hasher:
		return pbkdf2.NewPBKDF2SHA1Hasher().Verify(password, encoded)
	case PBKDF2SHA256Hasher:
		return pbkdf2.NewPBKDF2SHA256Hasher().Verify(password, encoded)
	}

	return false, ErrInvalidHasher
}

// MakePassword turns a plain-text password into a hash.
//
// If password is empty then return a concatenation
// of UnusablePasswordPrefix and a random string.
// If salt is empty then a randon string is generated.
// BCrypt algorithm ignores salt parameter.
// If hasher is "default" encode using default hasher.
func MakePassword(password, salt, hasher string) (string, error) {
	if password == "" {
		return UnusablePasswordPrefix + GetRandomString(UnusablePasswordSuffixLength), nil
	}

	if salt == "" {
		salt = GetRandomString(DefaultSaltSize)
	}

	if hasher == "default" {
		hasher = DefaultHasher
	}

	switch hasher {
	case Argon2Hasher:
		return argon2.NewArgon2Hasher().Encode(password, salt)
	case BCryptHasher:
		return bcrypt.NewBCryptHasher().Encode(password, salt)
	case BCryptSHA256Hasher:
		return bcrypt.NewBCryptSHA256Hasher().Encode(password, salt)
	case PBKDF2SHA1Hasher:
		return pbkdf2.NewPBKDF2SHA1Hasher().Encode(password, salt, 0)
	case PBKDF2SHA256Hasher:
		return pbkdf2.NewPBKDF2SHA256Hasher().Encode(password, salt, 0)
	}

	if IsValidHasher(hasher) {
		return "", ErrHasherNotImplemented
	}

	return "", ErrInvalidHasher
}
