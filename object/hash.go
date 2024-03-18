package object

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Hash [sha1.Size]byte

func HashFromString(s string) (Hash, error) {
	if len(s) != len(Hash{})*2 {
		return Hash{}, fmt.Errorf("invalid hash length: %d", len(s))
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, fmt.Errorf("decode string: %w", err)
	}

	var hash Hash
	copy(hash[:], b)

	return hash, nil
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
