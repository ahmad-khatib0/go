// Package ccm implements a CCM, Counter with CBC-MAC
// as per RFC 3610.
//
// See https://tools.ietf.org/html/rfc3610
//
// This code was lifted from https://github.com/bocajim/dtls/blob/a3300364a283fcb490d28a93d7fcfa7ba437fbbe/ccm/ccm.go
// and as such was not written by the Pions authors. Like Pions this
// code is licensed under MIT.
//
// A request for including CCM into the Go standard library
// can be found as issue #27484 on the https://github.com/golang/go/
// repository.
package ccm

import (
	"crypto/cipher"
	"errors"
)

// ccm represents a Counter with CBC-MAC with a specific key.
type ccm struct {
	b cipher.Block
	M uint8
	L uint8
}

const ccmBlockSize = 16

// CCM is a block cipher in Counter with CBC-MAC mode.
// Providing authenticated encryption with associated data via the cipher.AEAD interface.
type CCM interface {
	cipher.AEAD
	// MaxLength returns the maxium length of plaintext in calls to Seal.
	// The maximum length of ciphertext in calls to Open is MaxLength()+Overhead().
	// The maximum length is related to CCM's `L` parameter (15-noncesize) and
	// is 1<<(8*L) - 1 (but also limited by the maxium size of an int).
	MaxLength() int
}

var (
	errInvalidBlockSize = errors.New("ccm: NewCCM requires 128-bit block cipher")
	errInvalidTagSize   = errors.New("ccm: tagsize must be 4, 6, 8, 10, 12, 14, or 16")
	errInvalidNonceSize = errors.New("ccm: invalid nonce size")
)
