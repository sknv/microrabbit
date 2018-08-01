package xstrings

import (
	"math/rand"
	"time"
)

const (
	letters           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers           = "0123456789"
	lettersAndNumbers = letters + numbers
	charIdxBits       = 6                  // 6 bits to represent a letter index.
	charIdxMask       = 1<<charIdxBits - 1 // All 1-bits, as many as charIdxBits.
	charIdxMax        = 63 / charIdxBits   // # of letter indices fitting in 63 bits.
)

var (
	rndSrc = rand.NewSource(time.Now().UnixNano())
)

func RandomLetters(length int) string {
	return randomString(letters, length)
}

func RandomNumbers(length int) string {
	return randomString(numbers, length)
}

func RandomString(length int) string {
	return randomString(lettersAndNumbers, length)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func randomString(source string, length int) string {
	buffer := make([]byte, length)
	letterCount := len(source)
	// A rndSrc.Int63() generates 63 random bits, enough for charIdxMax characters.
	for i, cache, remain := length-1, rndSrc.Int63(), charIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rndSrc.Int63(), charIdxMax
		}
		if idx := int(cache & charIdxMask); idx < letterCount {
			buffer[i] = source[idx]
			i--
		}
		cache >>= charIdxBits
		remain--
	}
	return string(buffer)
}
