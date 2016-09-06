package interval_tree_test

import (
	"flag"
	"github.com/cockroachdb/cockroach/util/interval"
	"github.com/yaojingguo/benchmark-interval-tree/fixture"
	"testing"
)

var impl = flag.String("impl", "llrb", "interval tree implementation: llrb or btree")
var degree = flag.Int("degree", 32, "B-tree degree")

const (
	M = 1 << 20
)

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

func loadTree(b *testing.B, N int) (ivs []interval.Interface, tree *interval.LLRB) {
	ivs = fixture.GenN(b, N)
	tree = &interval.LLRB{Overlapper: interval.InclusiveOverlapper}
	for _, iv := range ivs {
		if err := tree.Insert(iv, false); err != nil {
			b.Fatalf("fast insert error: %s", err)
		}
	}
	tree.AdjustRanges()
	return
}

func loadRandomTree(b *testing.B, N int) (ivs []interval.Interface, tree *interval.LLRB) {
	ivs = fixture.RandomGenN(b, N)
	tree = &interval.LLRB{Overlapper: interval.InclusiveOverlapper}
	for _, iv := range ivs {
		if err := tree.Insert(iv, false); err != nil {
			b.Fatalf("fast insert error: %s", err)
		}
	}
	tree.AdjustRanges()
	return
}

func BenchmarkInsert(b *testing.B) {
	ivs := fixture.GenN(b, M)
	tree := NewTree()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, iv := range ivs {
			if err := tree.Insert(iv, false); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
	}
}

func BenchmarkFastInsert(b *testing.B) {
	ivs := fixture.GenN(b, M)
	tree := NewTree()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, iv := range ivs {
			if err := tree.Insert(iv, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
	}
}

func BenchmarkDelete(b *testing.B) {
	ivs, tree := loadTree(b, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, iv := range ivs {
			if err := tree.Delete(iv, false); err != nil {
				b.Fatalf("delete error: %s", err)
			}
		}
		b.StopTimer()
		if tree.Len() != 0 {
			b.Errorf("expectecd tree length %d, got %d", 0, tree.Len())
		}
		b.StartTimer()
	}
}

func BenchmarkGet(b *testing.B) {
	ivs, tree := loadTree(b, M)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, iv := range ivs {
			ptr := iv.(*fixture.Interval)
			tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
		}
	}
}

func BenchmarkNewTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_ = NewTree()
		}
	}
}

func BenchmarkRandomFixedInserts(b *testing.B) {
	ivs := fixture.RandomGenN(b, *fixture.Size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := NewTree()
		for _, e := range ivs {
			if err := tree.Insert(e, false); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
	}
}

func BenchmarkRandomFixedFastInserts(b *testing.B) {
	ivs := fixture.RandomGenN(b, *fixture.Size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := NewTree()
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
	}
}

func BenchmarkRandomFixedDeletes(b *testing.B) {
	ivs := fixture.RandomGenN(b, *fixture.Size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tree := NewTree()
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
		b.StartTimer()
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

func BenchmarkRandomFixedGets(b *testing.B) {
	ivs := fixture.RandomGenN(b, *fixture.Size)
	tree := NewTree()
	for _, e := range ivs {
		if err := tree.Insert(e, true); err != nil {
			b.Fatalf("insert error: %s", err)
		}
	}
	tree.AdjustRanges()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range ivs {
			ptr := e.(*fixture.Interval)
			tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
		}
	}
}

func BenchmarkRandomInsert(b *testing.B) {
	ivs := fixture.RandomGen(b)
	tree := NewTree()
	b.ResetTimer()
	for _, e := range ivs {
		if err := tree.Insert(e, false); err != nil {
			b.Fatalf("insert error: %s", err)
		}
	}
}

func BenchmarkRandomFastInsert(b *testing.B) {
	ivs := fixture.RandomGen(b)
	tree := NewTree()
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
		ptr := iv.(*fixture.Interval)
		tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
	}
}

func benchmarkFixedSizeGet(b *testing.B, N int) {
	ivs, tree := loadTree(b, N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, iv := range ivs {
			ptr := iv.(*fixture.Interval)
			tree.Get(interval.Range{ptr.R.Start, ptr.R.End})
		}
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
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x01}, End: interval.Comparable{0x02}}, uintptr(0)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x04}, End: interval.Comparable{0x06}}, uintptr(1)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x00}, End: interval.Comparable{0x02}}, uintptr(2)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x01}, End: interval.Comparable{0x06}}, uintptr(3)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x05}, End: interval.Comparable{0x15}}, uintptr(4)},
		&fixture.Interval{interval.Range{Start: interval.Comparable{0x25}, End: interval.Comparable{0x30}}, uintptr(5)},
	}
}

func BenchmarkInsertWithSmallTree(b *testing.B) {
	ivs := fewIntervals()
	for i := 0; i < b.N; i++ {
		tree := NewTree()
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
		tree := NewTree()
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
		tree := NewTree()
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
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
		tree := NewTree()
		for _, e := range ivs {
			if err := tree.Insert(e, true); err != nil {
				b.Fatalf("insert error: %s", err)
			}
		}
		tree.AdjustRanges()
		for _, e := range ivs {
			tree.Get(e.Range())
		}
	}
}
