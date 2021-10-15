package quantity

type Q struct {
	Number        float64
	UnitExponents UCombination

	// this is for display only, track
	// whether a value is input derived
	DerivedUnitsToUse UDerivedList
}
