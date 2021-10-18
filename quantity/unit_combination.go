package quantity

import (
	"bytes"
	"fmt"
	"sort"
)

type UCombination []UExp

func (u1 UCombination) HasOverlap(u2 UCombination) bool {
	unitUsed := make(map[int]bool)

	u1.Simplify()
	u2.Simplify()
	for _, u := range u1 {
		unitUsed[u.Unit.ID] = true
	}
	for _, u := range u2 {
		if _, ok := unitUsed[u.Unit.ID]; ok {
			return true
		}
	}
	return false
}

func (u UCombination) Inverse() {
	for i := range u {
		u[i].Exponent = -u[i].Exponent
	}
}

func (u UCombination) IsNoUnit() bool {
	u.Simplify()
	return len(u) == 0
}

func (u UCombination) Clone() UCombination {
	ret := make(UCombination, len(u))
	copy(ret, u)
	return ret
}

// Derives tries to express comb in exponents of u
func (u UCombination) Derive(comb UCombination) (remain UCombination, exp int) {

	maxNegativeExponent := 0
	minPositiveExponent := 0

	comb.Simplify()

	remain = comb.Clone()

	for _, derivedUnitElem := range u {
		derivedUnitID := derivedUnitElem.Unit.ID
		found := false
		for _, gotElem := range comb {
			if gotElem.Unit.ID == derivedUnitID {
				found = true
				exp := gotElem.Exponent / derivedUnitElem.Exponent
				if exp > 0 {
					if minPositiveExponent == 0 {
						minPositiveExponent = exp
					} else if exp < minPositiveExponent {
						minPositiveExponent = exp
					}
				} else if exp < 0 {
					if maxNegativeExponent == 0 {
						maxNegativeExponent = exp
					} else if exp > maxNegativeExponent {
						maxNegativeExponent = exp
					}
				} else {
					return remain, 0
				}
			}
		}
		if !found {
			return remain, 0
		}
	}

	if minPositiveExponent == 0 {
		exp = maxNegativeExponent
	} else if maxNegativeExponent == 0 {
		exp = minPositiveExponent
	} else {
		return remain, 0
	}

	for i, c := range remain {
		unitID := c.Unit.ID
		for _, uc := range u {
			if uc.Unit.ID == unitID {
				remain[i].Exponent -= uc.Exponent * exp
			}
		}
	}
	remain.Simplify()
	return
}

func (u *UCombination) Simplify() {
	tmp := make(map[U]int)
	for _, e := range *u {
		tmp[e.Unit] += e.Exponent
	}
	var res UCombination
	for unit, e := range tmp {
		if e != 0 {
			res = append(res, UExp{Unit: unit, Exponent: e})
		}
	}
	sort.Sort(res)
	*u = res
}

func (u *UCombination) Equal(u2 *UCombination) bool {
	u.Simplify()
	u2.Simplify()

	if len(*u) != len(*u2) {
		return false
	}

	for i := range *u {
		if (*u)[i] != (*u2)[i] {
			return false
		}
	}
	return true
}

// Len is the number of elements in the collection.
func (u UCombination) Len() int {
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
func (u UCombination) Less(i int, j int) bool {
	if u[i].Exponent == u[j].Exponent {
		return u[i].Unit.ID > u[j].Unit.ID
	}
	return u[i].Exponent > u[j].Exponent
}

// Swap swaps the elements with indexes i and j.
func (u UCombination) Swap(i int, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UCombination) String() string {
	u.Simplify()
	var b bytes.Buffer
	for i := range u {
		fmt.Fprintf(&b, "(%s)", u[i].Unit.Identifier)
		if u[i].Exponent != 1 {
			fmt.Fprintf(&b, "%d", u[i].Exponent)
		}
	}
	return b.String()
}
