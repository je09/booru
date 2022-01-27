package lib

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"math/big"
	"strings"
)

// AuthHash generates hash for authenticated users. To prevent parameter brute-forcing.
func AuthHash() string {
	// There is no reason for using 32 symbols long string. Just a number out of my head.
	s := randomString(32)
	hasher := sha512.New()
	hasher.Write([]byte(s))

	return hex.EncodeToString(hasher.Sum(nil))
}

// randomString generates a random string.
// We assume that function isn't public and size always stays in range [1, math.MaxInt32], so we don't test it.
func randomString(size int) string {
	symbols := []int32("abcdefghijklmonpqrstuvwABCDEFGHIJKLMNOPQRSTUVW1234567890!@#$%^&*()-=")
	var randString strings.Builder

	for i := 0; i < size; i++ {
		// Generate index of the symbol using cryptographically secure random.
		v, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		// We assume that number of characters in symbol string is lower than math.MaxInt32.
		r := int32(v.Int64())
		// Get random rune from range and write it down.
		randString.WriteRune(symbols[r])
	}

	return randString.String()
}
