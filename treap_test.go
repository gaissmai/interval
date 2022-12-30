package interval_test

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/gaissmai/interval"
)

type Int int

// implement interval.Interval
func (i Int) CompareUpper(j Int) int { return i.CompareLower(j) }

func (i Int) CompareLower(j Int) int {
	if i == j {
		return 0
	}
	if i < j {
		return -1
	}
	return 1
}

func TestNewTree(t *testing.T) {
	n := interval.NewTree[Int]()
	want := 0
	if n.Size() != want {
		t.Fatalf("NewTree(nil), got size: %d, expected: %d\n", n.Size(), want)
	}
	if n.Height() != want {
		t.Fatalf("NewTree(nil), got height: %d, expected: %d\n", n.Height(), want)
	}

	n = interval.NewTree([]Int{1, 2, 3}...)
	want = 3
	if n.Size() != want {
		t.Fatalf("NewTree(nil), got size: %d, expected: %d\n", n.Size(), want)
	}
}

func TestAscend(t *testing.T) {
	input := []Int{}
	for kk := 0; kk < 5000; kk++ {
		input = append(input, Int(rand.Int()))
	}

	n := interval.NewTree(input...)

	var inOrder []Int
	n.Ascend(func(n *interval.Tree[Int]) bool {
		inOrder = append(inOrder, n.Item())
		return true
	})
	if len(inOrder) != len(input) {
		t.Logf("%v", inOrder)
		t.Fatalf("Ascend: got %v, want %v", len(inOrder), len(input))
	}

	sort.Slice(input, func(i, j int) bool {
		return input[i].CompareLower(input[j]) <= 0
	})

	if !reflect.DeepEqual(input, inOrder) {
		t.Fatal("Ascend and sorted input diverged")
	}
}
