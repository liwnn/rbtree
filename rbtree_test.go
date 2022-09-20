package rbtree

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func init() {
	seed := time.Now().Unix()
	rand.Seed(seed)
}

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []Item) {
	out = make([]Item, 0, n)
	for _, v := range rand.Perm(n) {
		out = append(out, Int(v))
	}
	return
}

// rang returns an ordered list of Int items in the range [0, n).
func rang(n int) (out []Item) {
	for i := 0; i < n; i++ {
		out = append(out, Int(i))
	}
	return
}

func TestRBtree(t *testing.T) {
	tr := New()
	const treeSize = 10000
	for i := 0; i < 10; i++ {
		for _, item := range perm(treeSize) {
			tr.Insert(item)
		}
		if tr.Len() != treeSize {
			t.Fatal("insert failed", treeSize, tr.Len())
		}
		for _, item := range perm(treeSize) {
			if tr.Search(item) == nil {
				t.Fatal("has did not find item", item)
			}
		}
		for _, item := range perm(treeSize) {
			tr.Insert(item)
		}
		it := tr.NewAscendIterator()
		if min, want := it.Value(), Item(Int(0)); min != want {
			t.Fatalf("min: want %+v, got %+v", want, min)
		}

		for _, item := range perm(treeSize) {
			if tr.Delete(item) == nil {
				t.Fatalf("didn't find %v", item)
			}
		}
	}
}

func ExampleRBTree() {
	tr := New()
	for i := Int(0); i < 10; i++ {
		tr.Insert(i)
	}
	fmt.Println("len:       ", tr.Len())
	fmt.Println("search3:   ", tr.Search(Int(3)))
	fmt.Println("search100: ", tr.Search(Int(100)))
	fmt.Println("del4:      ", tr.Delete(Int(4)))
	fmt.Println("del100:    ", tr.Delete(Int(100)))
	tr.Insert(Int(5))
	tr.Insert(Int(100))
	fmt.Println("len:       ", tr.Len())
	fmt.Printf("for:        ")
	for it := tr.NewAscendIterator(); it.Valid(); it.Next() {
		fmt.Print(it.Value().(Int))
		fmt.Print(" ")
	}
	fmt.Println()
	// Output:
	// len:        10
	// search3:    3
	// search100:  <nil>
	// del4:       4
	// del100:     <nil>
	// len:        10
	// for:        0 1 2 3 5 6 7 8 9 100
}

func TestAscendIterator(t *testing.T) {
	tr := New()
	for _, v := range perm(100) {
		tr.Insert(v)
	}

	var got = make([]Item, 0, 100)
	for it := tr.NewAscendIterator(); it.Valid(); it.Next() {
		got = append(got, it.Value())
	}

	if want := rang(100); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

const benchmarkTreeSize = 10000

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		tr := New()
		for _, item := range insertP {
			tr.Insert(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	searchP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		tr := New()
		for _, v := range insertP {
			tr.Insert(v)
		}
		b.StartTimer()
		for _, item := range searchP {
			tr.Search(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkDeleteInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	tr := New()
	for _, item := range insertP {
		tr.Insert(item)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Delete(insertP[i%benchmarkTreeSize])
		tr.Insert(insertP[i%benchmarkTreeSize])
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	removeP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		tr := New()
		for _, v := range insertP {
			tr.Insert(v)
		}
		b.StartTimer()
		for _, item := range removeP {
			tr.Delete(item)
			i++
			if i >= b.N {
				return
			}
		}
		if tr.Len() > 0 {
			panic(tr.Len())
		}
	}
}

func BenchmarkAscend(b *testing.B) {
	tr := New()
	for _, item := range perm(benchmarkTreeSize) {
		tr.Insert(item)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Ascend(func(i Item) {

		})
	}
}

func BenchmarkAscendIterater(b *testing.B) {
	tr := New()
	for _, item := range perm(benchmarkTreeSize) {
		tr.Insert(item)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := Iterator{
			t: tr,
			x: tr.minimum(tr.root),
		}
		for ; it.Valid(); it.Next() {
		}
	}
}
