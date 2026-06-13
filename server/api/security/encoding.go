package security

import (
	"errors"

	"github.com/sqids/sqids-go"
)

type Encoder interface {
	Encode(int) (string, error)
	Decode(string) (int, error)
}

type squidEncoder struct {
	sq *sqids.Sqids
}

func NewSquidEncoder(sq *sqids.Sqids) Encoder {
	return squidEncoder{sq: sq}
}

func (encoder squidEncoder) Encode(id int) (string, error) {
	return encoder.sq.Encode([]uint64{uint64(id)})
}

func (encoder squidEncoder) Decode(hashId string) (int, error) {
	ids := encoder.sq.Decode(hashId)

	if len(ids) == 0 {
		return 0, errors.New("invalid id")
	}

	return int(ids[0]), nil
}
