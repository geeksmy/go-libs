package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/speps/go-hashids"
)

const (
	hashIDSalt      string = ""
	hashIDAlphabet  string = "hijklmnopqrst12345abcdefg67890uvwxyz"
	hashIDMinLength int    = 5
)

// shortID produces a " unique" 6 bytes long string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
func ShortID() string {
	b := make([]byte, 6)
	_, _ = io.ReadFull(rand.Reader, b)

	return base64.RawURLEncoding.EncodeToString(b)
}

func HashIDEncode(id int) string {
	hd := hashids.NewData()
	hd.Salt = hashIDSalt
	hd.Alphabet = hashIDAlphabet
	hd.MinLength = hashIDMinLength

	h, _ := hashids.NewWithData(hd)
	idStr, _ := h.Encode([]int{id})

	return idStr
}

func HashIDDecode(idStr string) (int, error) {
	hd := hashids.NewData()
	hd.Salt = hashIDSalt
	hd.Alphabet = hashIDAlphabet
	hd.MinLength = hashIDMinLength

	h, _ := hashids.NewWithData(hd)
	d, err := h.DecodeWithError(idStr)

	if err != nil {
		return 0, err
	}

	return d[0], nil
}
