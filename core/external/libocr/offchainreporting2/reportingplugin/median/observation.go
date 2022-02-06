package median

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
)

var i = big.NewInt

// Bounds on an ethereum int192
const byteWidth = 24
const bitWidth = byteWidth * 8

var MaxObservation = i(0).Sub(i(0).Lsh(i(1), bitWidth-1), i(1)) // 2**191 - 1
var MinObservation = i(0).Sub(i(0).Neg(MaxObservation), i(1))   // -2**191

func ToBytes(o *big.Int) ([]byte, error) {
	if o.Cmp(MaxObservation) > 0 || o.Cmp(MinObservation) < 0 {
		return nil, fmt.Errorf("value won't fit in int%v: 0x%x", bitWidth, o)
	}
	negative := o.Sign() < 0
	val := (&big.Int{})
	if negative {
		// compute two's complement as 2**192 - abs(o.v) = 2**192 + o.v
		val.SetInt64(1)
		val.Lsh(val, bitWidth)
		val.Add(val, o)
	} else {
		val.Set(o)
	}
	b := val.Bytes() // big-endian representation of abs(val)
	if len(b) > byteWidth {
		return nil, fmt.Errorf("b must fit in %v bytes", byteWidth)
	}
	b = bytes.Join([][]byte{bytes.Repeat([]byte{0}, byteWidth-len(b)), b}, []byte{})
	if len(b) != byteWidth {
		return nil, fmt.Errorf("wrong length; there must be an error in the padding of b: %v", b)
	}
	return b, nil
}

func ToBigInt(s []byte) (*big.Int, error) {
	if len(s) != byteWidth {
		return &big.Int{}, errors.Errorf("wrong length for serialized "+
			"Observation: length %d 0x%x", len(s), s)
	}
	val := (&big.Int{}).SetBytes(s)
	negative := val.Cmp(MaxObservation) > 0
	if negative {
		maxUint := (&big.Int{}).SetInt64(1)
		maxUint.Lsh(maxUint, bitWidth)
		val.Sub(maxUint, val)
		val.Neg(val)
	}
	return val, nil
}
