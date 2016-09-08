package interval_tree_test

import (
	"flag"
	"github.com/cockroachdb/cockroach/util/interval"
	"github.com/yaojingguo/benchmark-interval-tree/fixture"
	"testing"
)

var impl = flag.String("impl", "llrb", "interval tree implementation: llrb or btree")
var degree = flag.Int("degree", 32, "B-tree degree")
var size = flag.Int("size", 8, "tree size")

func NewTree() interval.Tree {
	switch *impl {
	case "llrb":
		return interval.Tree{Overlapper: interval.InclusiveOverlapper}
	// case "btree":
	//   return &interval.BTree{Overlapper: interval.InclusiveOverlapper, Degree: *degree}
	default:
		panic("invalid implementation")
	}
}

func NewLLRB() interval.Tree {
	return interval.Tree{Overlapper: interval.InclusiveOverlapper}
}

func BenchmarkLLRB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_ = interval.Tree{Overlapper: interval.InclusiveOverlapper}
			// _ = NewLLRB()
		}
	}
}

func loadTree(b *testing.B, ivs []interval.Interface) (tree interval.Tree) {
	tree = NewTree()
	for _, iv := range ivs {
		if err := tree.Insert(iv, false); err != nil {
			b.Fatalf("fast insert error: %s", err)
		}
	}
	tree.AdjustRanges()
	return
}

func rangeGroupRestIntervals() []interval.Interface {
	return []interval.Interface{
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x01}, End: interval.Comparable{0x02}}, uintptr(0)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x04}, End: interval.Comparable{0x06}}, uintptr(1)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x00}, End: interval.Comparable{0x02}}, uintptr(2)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x01}, End: interval.Comparable{0x06}}, uintptr(3)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x05}, End: interval.Comparable{0x15}}, uintptr(4)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x25}, End: interval.Comparable{0x30}}, uintptr(5)},
	}
}

// 1M
// 8
const (
	tiny  = 8
	large = 1 << 20
)

func runInsert(b *testing.B, ivs []interval.Interface) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := NewTree()
		for _, iv := range ivs {
			if err := tree.Insert(iv, false); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
	}
}

func runFastInsert(b *testing.B, ivs []interval.Interface) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := NewTree()
		for _, iv := range ivs {
			if err := tree.Insert(iv, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
	}
}

// If b.StopTimer and b.StartTimer are used to ignore the costs of inserts, this benchmark takes a
// too long time to finish.
func runDelete(b *testing.B, ivs []interval.Interface) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// b.StopTimer()
		tree := NewTree()
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
		// b.StartTimer()
		for _, iv := range ivs {
			if err := tree.Delete(iv, false); err != nil {
				b.Fatalf("delete error: %s", err)
			}
		}
		if tree.Len() != 0 {
			b.Errorf("expectecd tree length %d, got %d", 0, tree.Len())
		}
	}
}

func runGet(b *testing.B, ivs []interval.Interface) {
	tree := loadTree(b, ivs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, iv := range ivs {
			ptr := iv.(*fixture.Interval)
			tree.Get(ptr.R)
		}
	}
}

func benchmarkInsertN(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runInsert(b, ivs)
}

func benchmarkFastInsertN(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runFastInsert(b, ivs)
}

func benchmarkDeleteN(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runDelete(b, ivs)
}

func benchmarkGetN(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runGet(b, ivs)
}

func benchmarkRandomInsertN(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runInsert(b, ivs)
}

func benchmarkRandomFastInsertN(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runFastInsert(b, ivs)
}

func benchmarkRandomDeleteN(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runDelete(b, ivs)
}

func benchmarkRandomGetN(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runGet(b, ivs)
}

// Benchmarks

func BenchmarkNewTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_ = NewTree()
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	benchmarkInsertN(b, *size)
}

func BenchmarkFastInsert(b *testing.B) {
	benchmarkFastInsertN(b, *size)
}

func BenchmarkDelete(b *testing.B) {
	benchmarkDeleteN(b, *size)
}

func BenchmarkGet(b *testing.B) {
	benchmarkGetN(b, *size)
}

func BenchmarkRandomInsert(b *testing.B) {
	benchmarkRandomInsertN(b, *size)
}

func BenchmarkRandomFastInsert(b *testing.B) {
	benchmarkRandomFastInsertN(b, *size)
}

func BenchmarkRandomDelete(b *testing.B) {
	benchmarkRandomDeleteN(b, *size)
}

func BenchmarkRandomGet(b *testing.B) {
	benchmarkRandomGetN(b, *size)
}

func BenchmarkInsertWithRangeGroupTestIntervals(b *testing.B) {
	ivs := rangeGroupRestIntervals()
	runInsert(b, ivs)
}

func BenchmarkFastInsertWithRangeGroupTestIntervals(b *testing.B) {
	ivs := rangeGroupRestIntervals()
	runFastInsert(b, ivs)
}

func BenchmarkDeleteWithRangeGroupTestIntervals(b *testing.B) {
	ivs := rangeGroupRestIntervals()
	runDelete(b, ivs)
}

func BenchmarkGetWithRangeGroupTestIntervals(b *testing.B) {
	ivs := rangeGroupRestIntervals()
	runGet(b, ivs)
}

const (
	_4    = 4
	_8    = 8
	_100  = 100
	_1K   = 1000
	_10K  = 10 * _1K
	_100K = 10 * _10K
	_1M   = 10 * _100K
)

func BenchmarkInsert4(b *testing.B) {
	benchmarkInsertN(b, _4)
}

func BenchmarkFastInsert4(b *testing.B) {
	benchmarkFastInsertN(b, _4)
}

func BenchmarkDelete4(b *testing.B) {
	benchmarkDeleteN(b, _4)
}

func BenchmarkGet4(b *testing.B) {
	benchmarkGetN(b, _4)
}

func BenchmarkInsert8(b *testing.B) {
	benchmarkInsertN(b, _8)
}

func BenchmarkFastInsert8(b *testing.B) {
	benchmarkFastInsertN(b, _8)
}

func BenchmarkDelete8(b *testing.B) {
	benchmarkDeleteN(b, _8)
}

func BenchmarkGet8(b *testing.B) {
	benchmarkGetN(b, _8)
}

func BenchmarkInsert100(b *testing.B) {
	benchmarkInsertN(b, _100)
}

func BenchmarkFastInsert100(b *testing.B) {
	benchmarkFastInsertN(b, _100)
}

func BenchmarkDelete100(b *testing.B) {
	benchmarkDeleteN(b, _100)
}

func BenchmarkGet100(b *testing.B) {
	benchmarkGetN(b, _100)
}

func BenchmarkInsert1K(b *testing.B) {
	benchmarkInsertN(b, _1K)
}

func BenchmarkFastInsert1K(b *testing.B) {
	benchmarkFastInsertN(b, _1K)
}

func BenchmarkDelete1K(b *testing.B) {
	benchmarkDeleteN(b, _1K)
}

func BenchmarkGet1K(b *testing.B) {
	benchmarkGetN(b, _1K)
}

func BenchmarkInsert10K(b *testing.B) {
	benchmarkInsertN(b, _10K)
}

func BenchmarkFastInsert10K(b *testing.B) {
	benchmarkFastInsertN(b, _10K)
}

func BenchmarkDelete10K(b *testing.B) {
	benchmarkDeleteN(b, _10K)
}

func BenchmarkGet10K(b *testing.B) {
	benchmarkGetN(b, _10K)
}

func BenchmarkInsert100K(b *testing.B) {
	benchmarkInsertN(b, _100K)
}

func BenchmarkFastInsert100K(b *testing.B) {
	benchmarkFastInsertN(b, _100K)
}

func BenchmarkDelete100K(b *testing.B) {
	benchmarkDeleteN(b, _100K)
}

func BenchmarkGet100K(b *testing.B) {
	benchmarkGetN(b, _100K)
}

func BenchmarkInsert1M(b *testing.B) {
	benchmarkInsertN(b, _1M)
}

func BenchmarkFastInsert1M(b *testing.B) {
	benchmarkFastInsertN(b, _1M)
}

func BenchmarkDelete1M(b *testing.B) {
	benchmarkDeleteN(b, _1M)
}

func BenchmarkGet1M(b *testing.B) {
	benchmarkGetN(b, _1M)
}

// Random
func BenchmarkRandomInsert4(b *testing.B) {
	benchmarkRandomInsertN(b, _4)
}

func BenchmarkRandomFastInsert4(b *testing.B) {
	benchmarkRandomFastInsertN(b, _4)
}

func BenchmarkRandomDelete4(b *testing.B) {
	benchmarkRandomDeleteN(b, _4)
}

func BenchmarkRandomGet4(b *testing.B) {
	benchmarkRandomGetN(b, _4)
}

func BenchmarkRandomInsert8(b *testing.B) {
	benchmarkRandomInsertN(b, _8)
}

func BenchmarkRandomFastInsert8(b *testing.B) {
	benchmarkRandomFastInsertN(b, _8)
}

func BenchmarkRandomDelete8(b *testing.B) {
	benchmarkRandomDeleteN(b, _8)
}

func BenchmarkRandomGet8(b *testing.B) {
	benchmarkRandomGetN(b, _8)
}

func BenchmarkRandomInsert100(b *testing.B) {
	benchmarkRandomInsertN(b, _100)
}

func BenchmarkRandomFastInsert100(b *testing.B) {
	benchmarkRandomFastInsertN(b, _100)
}

func BenchmarkRandomDelete100(b *testing.B) {
	benchmarkRandomDeleteN(b, _100)
}

func BenchmarkRandomGet100(b *testing.B) {
	benchmarkRandomGetN(b, _100)
}

func BenchmarkRandomInsert1K(b *testing.B) {
	benchmarkRandomInsertN(b, _1K)
}

func BenchmarkRandomFastInsert1K(b *testing.B) {
	benchmarkRandomFastInsertN(b, _1K)
}

func BenchmarkRandomDelete1K(b *testing.B) {
	benchmarkRandomDeleteN(b, _1K)
}

func BenchmarkRandomGet1K(b *testing.B) {
	benchmarkRandomGetN(b, _1K)
}

func BenchmarkRandomInsert10K(b *testing.B) {
	benchmarkRandomInsertN(b, _10K)
}

func BenchmarkRandomFastInsert10K(b *testing.B) {
	benchmarkRandomFastInsertN(b, _10K)
}

func BenchmarkRandomDelete10K(b *testing.B) {
	benchmarkRandomDeleteN(b, _10K)
}

func BenchmarkRandomGet10K(b *testing.B) {
	benchmarkRandomGetN(b, _10K)
}

func BenchmarkRandomInsert100K(b *testing.B) {
	benchmarkRandomInsertN(b, _100K)
}

func BenchmarkRandomFastInsert100K(b *testing.B) {
	benchmarkRandomFastInsertN(b, _100K)
}

func BenchmarkRandomDelete100K(b *testing.B) {
	benchmarkRandomDeleteN(b, _100K)
}

/*
// Commented since it takes a too long time to finish.
func BenchmarkRandomGet100K(b *testing.B) {
	benchmarkRandomGetN(b, _100K)
}
*/

func BenchmarkRandomInsert1M(b *testing.B) {
	benchmarkRandomInsertN(b, _1M)
}

func BenchmarkRandomFastInsert1M(b *testing.B) {
	benchmarkRandomFastInsertN(b, _1M)
}

func BenchmarkRandomDelete1M(b *testing.B) {
	benchmarkRandomDeleteN(b, _1M)
}

/*
// Commented since it takes a too long time to finish.
func BenchmarkRandomGet1M(b *testing.B) {
	benchmarkRandomGetN(b, _1M)
}
*/
