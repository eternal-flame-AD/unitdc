package quantity

var (
	UnitGram = U{
		Identifier: "g",
		ID:         1,
	}
	UnitLiter = U{
		Identifier: "l",
		ID:         2,
	}
	UnitIU = U{
		Identifier: "iu",
		ID:         3,
	}
	UnitMeter = U{
		Identifier: "m",
		ID:         4,
	}
	UnitMole = U{
		Identifier: "mol",
		ID:         5,
	}
)

var (
	UnitDerivedGramEng = DeriveUnitWithEngineeringSymbolList(
		UnitGram,
		"m", "u", "n", "p",
	)
	UnitDerivedLiterEng = DeriveUnitWithEngineeringSymbolList(
		UnitLiter,
		"d", "m", "u", "n",
	)
	UnitDerivedMeterEng = DeriveUnitWithEngineeringSymbolList(
		UnitMeter,
		"c", "m", "u", "n",
	)
	UnitDerivedMoleEng = DeriveUnitWithEngineeringSymbolList(
		UnitMole,
		"m", "u", "n", "p",
	)
	UnitDerivedAmu = func() UDerived {
		res := UDerived{
			Offset:     0,
			Multiplier: 1,
			Identifier: "Da",
			UnitExponents: UCombination{
				{
					Unit:     UnitGram,
					Exponent: 1,
				},
				{
					Unit:     UnitMole,
					Exponent: -1,
				},
			},
		}
		res.UnitExponents.Simplify()
		return res
	}()
	UnitDerivedAmuEng = DeriveUnitWithEnginneringSymbolOnDerivedUnitList(
		UnitDerivedAmu,
		"k",
	)
	UnitDerivedMolar = func() UDerived {
		res := UDerived{
			Offset:     0,
			Multiplier: 1,
			Identifier: "M",
			UnitExponents: UCombination{
				{
					Unit:     UnitMole,
					Exponent: 1,
				},
				{
					Unit:     UnitLiter,
					Exponent: -1,
				},
			},
		}
		res.UnitExponents.Simplify()
		return res
	}()
	UnitDerivedMolarEng = DeriveUnitWithEnginneringSymbolOnDerivedUnitList(
		UnitDerivedMolar,
		"m", "u", "n", "p",
	)
)
