package util

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

func CreatePassword(password string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(999999) + 1
	saltstr := fmt.Sprintf("%06d", r)
	tmp := password + saltstr

	h := sha1.New()
	h.Write([]byte(tmp))
	bs := h.Sum(nil)
	ret := fmt.Sprintf("shal1$%s$%x", saltstr, bs)
	return ret, nil
}

func CreatePassWithRand(password string, rand string) (string, error) {
	tmp := password + rand
	h := sha1.New()
	h.Write([]byte(tmp))
	bs := h.Sum(nil)
	ret := fmt.Sprintf("shal1$%s$%x", rand, bs)
	return ret, nil
}
