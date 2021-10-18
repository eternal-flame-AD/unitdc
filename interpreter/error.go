package interpreter

import (
	"github.com/eternal-flame-ad/unitdc/localize"
	"github.com/eternal-flame-ad/unitdc/quantity"
	"github.com/eternal-flame-ad/unitdc/syntax"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ErrEmptyStack struct {
}

func (e ErrEmptyStack) Error() string {
	return localize.Localizer().MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "InterpreterError_StackEmpty",
			Other: "Stack Empty",
		},
	})
}

type ErrIncompatibleUnit struct {
	TargetUnit    quantity.UCombination
	OffendingUnit quantity.UCombination
}

func (e ErrIncompatibleUnit) Error() string {
	if e.TargetUnit != nil {
		return localize.Localizer().MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "InterpreterError_IncompatibleUnitConvert",
				Other: "incompatible units: could not coerce {{.OffendingUnit}} to {{.TargetUnit}}",
			},
			TemplateData: map[string]interface{}{
				"OffendingUnit": e.OffendingUnit.String(),
				"TargetUnit":    e.TargetUnit.String(),
			},
		})
	}
	return localize.Localizer().MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "InterpreterError_IncompatibleUnitUnacceptable",
			Other: "incompatible units: {{.OffendingUnit}} is unacceptable for this operation",
		},
		TemplateData: map[string]interface{}{
			"OffendingUnit": e.OffendingUnit.String(),
		},
	})
}

type ErrUnknownOperation struct {
	Token syntax.Token
}

func (e ErrUnknownOperation) Error() string {
	return localize.Localizer().MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "InterpreterError_UnknownOperation",
			Other: "undefined operation: {{.Operation}}",
		},
		TemplateData: map[string]interface{}{
			"Operation": e.Token.String(),
		},
	})
}

type ErrUnknownUnit struct {
	UnitIdentifier string
}

func (e ErrUnknownUnit) Error() string {
	return localize.Localizer().MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "InterpreterError_UnknownUnit",
			Other: "undefined unit: {{.UnitIdent}}",
		},
		TemplateData: map[string]interface{}{
			"UnitIdent": e.UnitIdentifier,
		},
	})
}
