package period

import "fmt"

// example interval
type Ival [2]int

// implement interval.Interface
func (p Ival) CompareLower(q Ival) int { return cmp(p[0], q[0]) }
func (p Ival) CompareUpper(q Ival) int { return cmp(p[1], q[1]) }

// fmt.Stringer for formattting, not required for interval.Interface
func (p Ival) String() string {
	return fmt.Sprintf("%d...%d", p[0], p[1])
}

// little helper, compare two ints
func cmp(a, b int) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return -1
	}
	return 1
}
