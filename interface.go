package interval

// Interface is the type constraint for generic interval items.
//
// There exists thirteen basic relations between any two intervals in one dimension,
// see [Allen's Interval Algebra].
//
//  =================================================================|
//  |  visualization         | ll | rr | lr | rl | description       |
//  =================================================================|
//  |  A1---A2               | -1 | -1 | -1 | -1 | A precedes B      |
//  |           B1---B2      |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1---A2               | -1 | -1 | -1 |  0 | A meets B         |
//  |       B1---B2          |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1----A2              | -1 | -1 | -1 |  1 | A overlaps B      |
//  |     B1-----B2          |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1------A2            | -1 |  0 | -1 |  1 | A finished by B   |
//  |     B1---B2            |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1--------A2          | -1 |  1 | -1 |  1 | A contains B      |
//  |     B1---B2            |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1---A2               |  0 | -1 | -1 |  1 | A starts B        |
//  |  B1-------B2           |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1----A2              |  0 |  0 | -1 |  1 | A equals B        |
//  |  B1----B2              |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1--------A2          |  0 |  1 | -1 |  1 | A started by B    |
//  |  B1----B2              |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |    A1--A2              |  1 | -1 | -1 |  1 | A during B        |
//  |  B1--------B2          |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |      A1----A2          |  1 |  0 | -1 |  1 | A finishes B      |
//  |  B1--------B2          |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |      A1----A2          |  1 |  1 | -1 |  1 | A overlapped by B |
//  |  B1-----B2             |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |       A1----A2         |  1 |  1 |  0 |  1 | A met by B        |
//  |  B1---B2               |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |          A1---A2       |  1 |  1 |  1 |  1 | A preceded by B   |
//  |  B1---B2               |    |    |    |    |                   |
//  -------------------------|----------------------------------------
//
// [Allen's Interval Algebra]: https://www.ics.uci.edu/~alspaugh/cls/shr/allen.html
//
type Interface[T any] interface {
	// Compare the left (l) and right (r) points of two intervals and returns four integers with values (-1, 0, +1).
	//
	// The result will be one of the thirteen possibilities in the interval relation table.
	//
	//  ll: left  point interval A compared with left  point interval B (-1, 0, +1)
	//  rr: right point interval A compared with right point interval B (-1, 0, +1)
	//  lr: left  point interval A compared with right point interval B (-1, 0, +1)
	//  rl: right point interval A compared with left  point interval B (-1, 0, +1)
	//
	Compare(T) (ll, rr, lr, rl int)
}

// compare is for sorting keys into the BST, the sort key is the left point of the intervals.
// If the left point is equal, sort the supersets to the left (definite order).
//
//  e.g. special treatment for all relations with ll == 0
//  =================================================================|
//  |  visualization         | ll | rr | lr | rl | description       |
//  =================================================================|
//  ~                                                                ~
//  -------------------------|---------------------------------------|
//  |  A1---A2               |  0 | -1 | -1 |  1 | A starts B        | => A is covered by B, => sort A to the right
//  |  B1-------B2           |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1----A2              |  0 |  0 | -1 |  1 | A equals B        | => equality
//  |  B1----B2              |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1--------A2          |  0 |  1 | -1 |  1 | A started by B    | => A covers B, => sort A to the left
//  |  B1----B2              |    |    |    |    |                   |
//  -------------------------|----------------------------------------
//
func compare[T Interface[T]](a, b T) int {
	ll, rr, _, _ := a.Compare(b)
	switch {
	case ll == 0:
		// identical left point, sort supersets to the left, make a definite order
		// if rr == 0, -rr makes no difference
		return -rr
	default:
		// interval left point is the BST sort key
		return ll
	}
}

// covers, returns true if a covers b.
//
//  =================================================================|
//  |  visualization         | ll | rr | lr | rl | description       |
//  =================================================================|
//  ~                                                                ~
//  -------------------------|---------------------------------------|
//  |  A1------A2            | -1 |  0 | -1 |  1 | A finished by B   |
//  |     B1---B2            |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1--------A2          | -1 |  1 | -1 |  1 | A contains by B   |
//  |     B1---B2            |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1--------A2          |  0 |  1 | -1 |  1 | A started by B    |
//  |  B1----B2              |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |  A1----A2              |  0 |  0 | -1 |  1 | A equals B        |
//  |  B1----B2              |    |    |    |    |                   |
//  -------------------------|----------------------------------------
//
func covers[T Interface[T]](a, b T) bool {
	ll, rr, _, _ := a.Compare(b)
	return ll <= 0 && rr >= 0
}

// intersects, returns false if the intervals does precede each other.
//
//  =================================================================|
//  |  visualization         | ll | rr | lr | rl | description       |
//  =================================================================|
//  ~                                                                ~
//  -------------------------|---------------------------------------|
//  |  A1---A2               | -1 | -1 | -1 | -1 | A precedes B      | => false
//  |           B1---B2      |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//  |          A1---A2       |  1 |  1 |  1 |  1 | A preceded by B   | => false
//  |  B1---B2               |    |    |    |    |                   |
//  -------------------------|---------------------------------------|
//
func intersects[T Interface[T]](a, b T) bool {
	ll, rr, lr, rl := a.Compare(b)
	return !((ll == -1 && rr == -1 && lr == -1 && rl == -1) || (ll == 1 && rr == 1 && lr == 1 && rl == 1))
}

// cmpRR, compares just the right point of the intervals.
func cmpRR[T Interface[T]](a, b T) int {
	_, rr, _, _ := a.Compare(b)
	return rr
}

// cmpLR, compares just the left point from a with right point from b.
func cmpLR[T Interface[T]](a, b T) int {
	_, _, lr, _ := a.Compare(b)
	return lr
}
