package quantity

import "sort"

type UnitDisplay struct {
	Identifier string
	Exponent   int
}

type UnitDisplayList []UnitDisplay

// Len is the number of elements in the collection.
func (u UnitDisplayList) Len() int {
	return len(u)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (u UnitDisplayList) Less(i int, j int) bool {
	if u[i].Exponent == u[j].Exponent {
		return u[i].Identifier > u[j].Identifier
	}
	return u[i].Exponent > u[j].Exponent
}

// Swap swaps the elements with indexes i and j.
func (u UnitDisplayList) Swap(i int, j int) {
	u[i], u[j] = u[j], u[i]
}

func (q Q) Format() (num float64, res UnitDisplayList) {
	num = q.Number
	comb := q.UnitExponents

	for _, d := range q.DerivedUnitsToUse {
		remain, exp := d.UnitExponents.Derive(comb)
		if exp != 0 {
			comb = remain

			for i := 0; i < exp; i++ {
				num /= d.Multiplier
				num -= d.Offset
			}
			for i := 0; i > exp; i-- {

				num += d.Offset
				num *= d.Multiplier
			}
			res = append(res, UnitDisplay{
				Identifier: d.Identifier,
				Exponent:   exp,
			})
		}
	}
	comb.Simplify()
	for _, u := range comb {
		res = append(res, UnitDisplay{
			Identifier: u.Unit.Identifier,
			Exponent:   u.Exponent,
		})
	}
	sort.Sort(res)
	return
}
