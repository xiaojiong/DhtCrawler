package godht

import (
	"crypto/sha1"
	"encoding/hex"

	"io"
	"math/big"
	"math/rand"
	"time"
)

type Id []byte

func (id Id) String() string {
	return hex.EncodeToString(id)
}

func (id Id) Int() *big.Int {
	return big.NewInt(0).SetBytes(id)
}

func (id Id) Neighbor() Id {
	randId := []byte(GenerateID())
	byteId := []byte(id)
	return append(byteId[0:12], randId[0:8]...)
}
func GenerateID() Id {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	hash := sha1.New()
	io.WriteString(hash, time.Now().String())
	io.WriteString(hash, string(random.Int()))

	return hash.Sum(nil)
}
