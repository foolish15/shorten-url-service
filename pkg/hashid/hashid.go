package hashid

import (
	"github.com/speps/go-hashids"
)

func getHD() (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Salt = "hashid pkg salt"
	hd.MinLength = 16
	return hashids.NewWithData(hd)
}

func Encrypt(input ...int) (string, error) {
	h, err := getHD()
	if err != nil {
		return "", err
	}
	return h.Encode(input)
}

func Decrypt(input string) ([]int, error) {
	h, err := getHD()
	if err != nil {
		return nil, err
	}

	o := h.Decode(input)
	return o, nil
}
