package quantity

import (
	"fmt"
	"sort"
)

type UDerived struct {
	Identifier    string
	Offset        float64
	Multiplier    float64
	UnitExponents UCombination
}

func (u UDerived) Simplifies(comb UCombination) {

}

func (u1 UDerived) Compatible(u2 UDerived) (ret bool) {
	for _, baseu := range u1.UnitExponents {
		for _, baseu2 := range u2.UnitExponents {
			if baseu.Unit.ID == baseu2.Unit.ID {
				return false
			}
		}
	}
	return true
}

type UDerivedList []UDerived

func (u UDerivedList) Compatible(u2 UDerived) bool {
	for _, ud := range u {
		if !ud.Compatible(u2) {
			return false
		}
	}
	return true
}

func (u UDerivedList) Clone() UDerivedList {
	res := make(UDerivedList, len(u))
	copy(res, u)
	return res
}

func CombineUDerivedLists(list ...UDerivedList) UDerivedList {
	res := list[0].Clone()
	for i := 1; i < len(list); i++ {
		for _, ud := range list[i] {
			if res.Compatible(ud) {
				res = append(res, ud)
			}
		}
	}
	return res
}

// Len is the number of elements in the collection.
func (u UDerivedList) Len() int {
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
func (u UDerivedList) Less(i int, j int) bool {
	return u[i].Identifier < u[j].Identifier
}

// Swap swaps the elements with indexes i and j.
func (u UDerivedList) Swap(i int, j int) {
	u[i], u[j] = u[j], u[i]
}

// IsComplex returns if the derived unit consist of more
// than one base unit
func (u UDerived) IsComplex() bool {
	return len(u.UnitExponents) > 1
}

func EqualUDerivedLists(l1 UDerivedList, l2 UDerivedList) bool {
	sort.Sort(l1)
	sort.Sort(l2)
	if len(l1) != len(l2) {
		return false
	}
	for i := range l1 {
		if l1[i].Identifier != l2[i].Identifier {
			return false
		}
	}
	return true
}

func DeriveUnitWithEngineeringSymbol(symbol string, base U) UDerived {
	var multiplier float64
	switch symbol {
	case "k":
		multiplier = 1e3
	case "d":
		multiplier = 1e-1
	case "c":
		multiplier = 1e-2
	case "m":
		multiplier = 1e-3
	case "u":
		multiplier = 1e-6
	case "n":
		multiplier = 1e-9
	case "p":
		multiplier = 1e-12
	default:
		panic("unknown engineering symbol")
	}
	return UDerived{
		Identifier:    symbol + base.Identifier,
		Offset:        0,
		Multiplier:    multiplier,
		UnitExponents: UCombination{UExp{Unit: base, Exponent: 1}},
	}
}

func DeriveUnitWithEngineeringSymbolList(base U, symbols ...string) UDerivedList {
	var res UDerivedList
	for _, s := range symbols {
		res = append(res, DeriveUnitWithEngineeringSymbol(s, base))
	}
	return res
}

func DeriveUnitWithEnginneringSymbolOnDerivedUnit(symbol string, base UDerived) UDerived {
	var multiplier float64
	switch symbol {
	case "k":
		multiplier = 1e3
	case "d":
		multiplier = 1e-1
	case "c":
		multiplier = 1e-2
	case "m":
		multiplier = 1e-3
	case "u":
		multiplier = 1e-6
	case "n":
		multiplier = 1e-9
	case "p":
		multiplier = 1e-12
	default:
		panic(fmt.Sprintf("unknown engineering symbol %s", symbol))
	}
	return UDerived{
		Identifier:    symbol + base.Identifier,
		Offset:        base.Offset,
		Multiplier:    base.Multiplier * multiplier,
		UnitExponents: base.UnitExponents.Clone(),
	}

}

func DeriveUnitWithEnginneringSymbolOnDerivedUnitList(base UDerived, symbols ...string) UDerivedList {
	var res UDerivedList
	for _, s := range symbols {
		res = append(res, DeriveUnitWithEnginneringSymbolOnDerivedUnit(s, base))
	}
	return res
}
