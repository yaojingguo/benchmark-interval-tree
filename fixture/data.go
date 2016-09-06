package fixture

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/cockroachdb/cockroach/util/interval"
	"math/rand"
	"testing"
	// "time"
)

var length = flag.Int("length", 1024, "max byte slice length")
var Size = flag.Int("size", 8, "tree size")

const (
	intervalLen = 10
)

type Interval struct {
	R   interval.Range
	ID_ uintptr
}

func init() {
	// rand.Seed(time.Now().UnixNano())
}

func (iv *Interval) Range() interval.Range {
	return iv.R
}

func (iv *Interval) ID() uintptr {
	return iv.ID_
}

func (iv *Interval) String() string {
	return fmt.Sprintf("%v-%d", iv.Range(), iv.ID())
}

func ToBytes(b *testing.B, n uint32) interval.Comparable {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, n); err != nil {
		b.Fatalf("binary.Write error: %s", err)
	}
	return interval.Comparable(buf.Bytes())
}

func createInterval(b *testing.B, start, end uint32) interval.Interface {
	iv := &Interval{interval.Range{ToBytes(b, start), ToBytes(b, end)}, uintptr(start)}
	return iv
}

func Gen(b *testing.B) []interval.Interface {
	return GenN(b, b.N)
}

func GenN(b *testing.B, N int) (ivs []interval.Interface) {
	for i := 0; i < N; i++ {
		n := uint32(i)
		ivs = append(ivs, createInterval(b, n, n+intervalLen))
	}
	return
}

func getRandomRange(b *testing.B, n int) interval.Range {
	s1 := getRandomByteSlice(b, n)
	s2 := getRandomByteSlice(b, n)
	cmp := bytes.Compare(s1, s2)
	for cmp == 0 {
		s2 = getRandomByteSlice(b, n)
		cmp = bytes.Compare(s1, s2)
	}
	if cmp < 0 {
		return interval.Range{Start: s1, End: s2}
	}
	return interval.Range{Start: s2, End: s1}
}

func getRandomByteSlice(b *testing.B, n int) interval.Comparable {
	length := rand.Intn(n) + 1
	s := make(interval.Comparable, length)
	_, err := rand.Read(s)
	if err != nil {
		b.Fatalf("could not create random byte slice: %v", err)
	}
	return s
}

func RandomGenN(b *testing.B, N int) (ivs []interval.Interface) {
	for i := 0; i < N; i++ {
		iv := &Interval{getRandomRange(b, *length), uintptr(i)}
		ivs = append(ivs, iv)
	}
	return
}

func RandomGen(b *testing.B) (ivs []interval.Interface) {
	return RandomGenN(b, b.N)
}
