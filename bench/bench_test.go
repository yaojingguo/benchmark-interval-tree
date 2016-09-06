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
		return &interval.LLRB{Overlapper: interval.InclusiveOverlapper}
	case "btree":
		return interval.NewBTreeWithDegree(interval.InclusiveOverlapper, *degree)
	default:
		panic("invalid implementation")
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

// If b.StopTimer and b.StartTimer are used to ignore the costs of inserts. "go test -bench ." takes
// a long time to run.
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
			tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
		}
	}
}

func benchmarkInsertN(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runInsert(b, ivs)
}

func benchmarkFastInsert(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runFastInsert(b, ivs)
}

func benchmarkDelete(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runDelete(b, ivs)
}

func benchmarkGet(b *testing.B, N int) {
	ivs := fixture.GenN(b, N)
	runGet(b, ivs)
}

func benchmarkRandomInsertN(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runInsert(b, ivs)
}

func benchmarkRandomFastInsert(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runFastInsert(b, ivs)
}

func benchmarkRandomDelete(b *testing.B, N int) {
	ivs := fixture.RandomGenN(b, N)
	runDelete(b, ivs)
}

func benchmarkRandomGet(b *testing.B, N int) {
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
	benchmarkFastInsert(b, *size)
}

func BenchmarkDelete(b *testing.B) {
	benchmarkDelete(b, *size)
}

func BenchmarkGet(b *testing.B) {
	benchmarkGet(b, *size)
}

func BenchmarkRandomInsert(b *testing.B) {
	benchmarkRandomInsertN(b, *size)
}

func BenchmarkRandomFastInsert(b *testing.B) {
	benchmarkRandomFastInsert(b, *size)
}

func BenchmarkRandomDelete(b *testing.B) {
	benchmarkRandomDelete(b, *size)
}

func BenchmarkRandomGet(b *testing.B) {
	benchmarkRandomGet(b, *size)
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
	_Tiny = 8
	_K    = 1000
	_10K  = 10 * _K
	_100K = 10 * _10K
	_1M   = 10 * _100K
)
