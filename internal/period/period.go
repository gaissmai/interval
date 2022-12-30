package period

import "fmt"

// example interval
type Ival struct {
	Start int
	Stop  int
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

// implement interval.Interface
func (p Ival) CompareLower(q Ival) int { return cmp(p.Start, q.Start) }
func (p Ival) CompareUpper(q Ival) int { return cmp(p.Stop, q.Stop) }

// fmt.Stringer for formattting, not required for interval.Interface
func (p Ival) String() string {
	return fmt.Sprintf("%d...%d", p.Start, p.Stop)
}
