package types

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/snowflake"
	"golang.org/x/crypto/xtea"
)

// UidGenerator holds snowflake and encryption paramenets.
type UidGenerator struct {
	seq    *snowflake.SnowFlake
	cipher *xtea.Cipher
}

// Init initialises the Uid generator
func NewUID(workerID uint, key []byte) (*UidGenerator, error) {
	var err error
	var ug UidGenerator

	ug.seq, err = snowflake.NewSnowFlake(uint32(workerID))
	ug.cipher, err = xtea.NewCipher(key)

	return &ug, err
}

// Get generates a unique weakly-encryped random-looking ID.
//
// # The Uid is a unit64 with the highest bit possibly set which makes it
//
// incompatible with go's pre-1.9 sql package.
func (ug *UidGenerator) Get() Uid {
	buf, err := getIDBuffer(ug)
	if err != nil {
		return ZeroUid
	}

	return Uid(binary.LittleEndian.Uint64(buf))
}

// GetStr generates the same unique ID as Get then returns it as
// base64-encoded string. Slightly more efficient than calling Get()
// then base64-encoding the result.
func (ug *UidGenerator) GetStr() string {
	buf, err := getIDBuffer(ug)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
}

// DecodeUid takes an encrypted Uid and decrypts it into a non-negative int64.
//
// This is needed for go/sql compatibility where uint64 with high bit
//
// set is unsupported and possibly for other uses such as MySQL's recommendation
//
// for sequential primary keys.
func (ug *UidGenerator) DecodeUid(uid Uid) int64 {
	if uid.IsZero() {
		return 0
	}

	src := make([]byte, 8)
	dst := make([]byte, 8)
	binary.LittleEndian.PutUint64(src, uint64(uid))
	ug.cipher.Decrypt(dst, src)
	return int64(binary.LittleEndian.Uint64(dst))
}

// EncodeUid takes a positive int64 and encrypts it into a Uid.
//
// This is needed for go/sql compatibility where uint64 with high bit
//
// set is unsupported  and possibly for other uses such as MySQL's recommendation
//
// for sequential primary keys.
func (ug *UidGenerator) EncodeUid(val int64) Uid {
	if val == 0 {
		return ZeroUid
	}

	src := make([]byte, 8)
	dst := make([]byte, 8)
	binary.LittleEndian.PutUint64(src, uint64(val))
	ug.cipher.Encrypt(dst, src)
	return Uid(binary.LittleEndian.Uint64(dst))
}

// getIdBuffer returns a byte array holding the Uid bytes
func getIDBuffer(ug *UidGenerator) ([]byte, error) {
	var id uint64
	var err error

	if id, err = ug.seq.Next(); err != nil {
		return nil, err
	}

	src := make([]byte, 8)
	dst := make([]byte, 8)

	// INFO: Endianness means that the bytes in computer memory are read in a certain order.
	binary.LittleEndian.PutUint64(src, id)
	ug.cipher.Encrypt(dst, src)

	return dst, nil
}

func (ug *UidGenerator) GetUidString() string {
	return ug.GetStr()
}
