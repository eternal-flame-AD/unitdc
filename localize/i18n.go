package localize

import (
	_ "embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed active.ja.json
var i18nJaJSON []byte

var localizer *i18n.Localizer

func Localizer() *i18n.Localizer {
	return localizer
}

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	//bundle.MustParseMessageFileBytes(i18nEnJSON, "active.en.json")
	bundle.MustParseMessageFileBytes(i18nJaJSON, "active.ja.json")
	localizer = i18n.NewLocalizer(bundle, getLang())
}
