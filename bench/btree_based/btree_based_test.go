package btree_based

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/cockroachdb/cockroach/util/interval"
	"math/rand"
	"testing"
	"time"
)

const (
	intervalLen = 10
)

type Interval struct {
	R  interval.Range
	id uintptr
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (iv *Interval) Range() interval.Range {
	return iv.R
}

func (iv *Interval) ID() uintptr {
	return iv.id
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

var degree = flag.Int("degree", 32, "B-tree degree")
var length = flag.Int("length", 1024, "max byte slice length")

func loadTree(b *testing.B, N int) (ivs []interval.Interface, tree *interval.BTree) {
	ivs = GenN(b, N)
	tree = interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	for _, iv := range ivs {
		if err := tree.Insert(iv, false); err != nil {
			b.Fatalf("fast insert error: %s", err)
		}
	}
	tree.AdjustRanges()
	return
}

func loadRandomTree(b *testing.B, N int) (ivs []interval.Interface, tree *interval.BTree) {
	ivs = RandomGenN(b, N)
	tree = interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	for _, iv := range ivs {
		if err := tree.Insert(iv, false); err != nil {
			b.Fatalf("fast insert error: %s", err)
		}
	}
	tree.AdjustRanges()
	return
}
func BenchmarkInsert(b *testing.B) {
	ivs := Gen(b)
	tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	b.ResetTimer()
	for _, e := range ivs {
		if err := tree.Insert(e, false); err != nil {
			b.Fatalf("insert error: %s", err)
		}
	}
}

func BenchmarkFastInsert(b *testing.B) {
	ivs := Gen(b)
	tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	b.ResetTimer()
	for _, iv := range ivs {
		if err := tree.Insert(iv, true); err != nil {
			b.Fatalf("insert error: %s", err)
		}
	}
	tree.AdjustRanges()
}

func BenchmarkDelete(b *testing.B) {
	ivs, tree := loadTree(b, b.N)
	b.ResetTimer()
	for _, iv := range ivs {
		if err := tree.Delete(iv, false); err != nil {
			b.Fatalf("delete error: %s", err)
		}
	}
	if tree.Len() != 0 {
		b.Errorf("expectecd tree length %d, got %d", 0, tree.Len())
	}
}

func BenchmarkGet(b *testing.B) {
	ivs, tree := loadTree(b, b.N)
	b.ResetTimer()
	for _, iv := range ivs {
		ptr := iv.(*Interval)
		tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
	}
}

func BenchmarkRandomInsert(b *testing.B) {
	ivs := RandomGen(b)
	tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	b.ResetTimer()
	for _, e := range ivs {
		if err := tree.Insert(e, false); err != nil {
			b.Fatalf("insert error: %s", err)
		}
	}
}

func BenchmarkRandomFastInsert(b *testing.B) {
	ivs := RandomGen(b)
	tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	b.ResetTimer()
	for _, iv := range ivs {
		if err := tree.Insert(iv, true); err != nil {
			b.Fatalf("insert error: %s", err)
		}
	}
	tree.AdjustRanges()
}

func BenchmarkRandomDelete(b *testing.B) {
	ivs, tree := loadRandomTree(b, b.N)
	b.ResetTimer()
	for _, iv := range ivs {
		if err := tree.Delete(iv, false); err != nil {
			b.Fatalf("delete error: %s", err)
		}
	}
	if tree.Len() != 0 {
		b.Errorf("expectecd tree length %d, got %d", 0, tree.Len())
	}
}

func BenchmarkRandomGet(b *testing.B) {
	ivs, tree := loadRandomTree(b, b.N)
	b.ResetTimer()
	for _, iv := range ivs {
		ptr := iv.(*Interval)
		tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
	}
}

func benchmarkFixedSizeGet(b *testing.B, N int) {
	ivs, tree := loadTree(b, N)
	b.ResetTimer()
	iLen := len(ivs)
	for i := 0; i < b.N; i++ {
		iv := ivs[i%iLen]
		ptr := iv.(*Interval)
		tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
	}
}

func BenchmarkGetFrom1k(b *testing.B) {
	benchmarkFixedSizeGet(b, 1*1000)
}

func BenchmarkGetFrom10k(b *testing.B) {
	benchmarkFixedSizeGet(b, 10*1000)
}

func BenchmarkGetFrom100k(b *testing.B) {
	benchmarkFixedSizeGet(b, 100*1000)
}

func BenchmarkGetFrom1000k(b *testing.B) {
	benchmarkFixedSizeGet(b, 1000*1000)
}

func fewIntervals() []interval.Interface {
	return []interval.Interface{
		&Interval{interval.Range{Start: interval.Comparable{0x01}, End: interval.Comparable{0x02}}, uintptr(0)},
		&Interval{interval.Range{Start: interval.Comparable{0x04}, End: interval.Comparable{0x06}}, uintptr(1)},
		&Interval{interval.Range{Start: interval.Comparable{0x00}, End: interval.Comparable{0x02}}, uintptr(2)},
		&Interval{interval.Range{Start: interval.Comparable{0x01}, End: interval.Comparable{0x06}}, uintptr(3)},
		&Interval{interval.Range{Start: interval.Comparable{0x05}, End: interval.Comparable{0x15}}, uintptr(4)},
		&Interval{interval.Range{Start: interval.Comparable{0x25}, End: interval.Comparable{0x30}}, uintptr(5)},
	}
}

func BenchmarkInsertWithSmallTree(b *testing.B) {
	ivs := fewIntervals()
	for i := 0; i < b.N; i++ {
		tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
		for _, e := range ivs {
			if err := tree.Insert(e, false); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
	}
}

func BenchmarkFastInsertWithSmallTree(b *testing.B) {
	ivs := fewIntervals()
	for i := 0; i < b.N; i++ {
		tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
	}
}

// If b.StopTimer and b.StartTimer are used to ignore the costs of inserts. "go test -bench ." takes
// a long time to run.
func BenchmarkInsertAndDeleteWithSmallTree(b *testing.B) {
	ivs := fewIntervals()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
		b.StartTimer()
		for _, e := range ivs {
			if err := tree.Delete(e, false); err != nil {
				b.Fatalf("delete error: %s", err)
			}
		}
		if tree.Len() != 0 {
			b.Errorf("expectecd tree length %d, got %d", 0, tree.Len())
		}
	}
}

// If b.StopTimer and b.StartTimer are used to ignore the costs of inserts. "go test -bench ." takes
// a long time to run.
func BenchmarkInsertAndGetWithSmallTree(b *testing.B) {
	ivs := fewIntervals()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tree := interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
		b.StartTimer()
		for _, e := range ivs {
			tree.Get(e.Range())
		}
	}
}
