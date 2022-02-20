package helper

import (
	"github.com/speps/go-hashids/v2"
)

const Salt = "StoryServiceHashSalt"

func GenerateHashId() string {
	hd := hashids.NewData()
	hd.Salt = Salt
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{45, 434, 1313, 99})

	return e
}
