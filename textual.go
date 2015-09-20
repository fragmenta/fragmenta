package main

import (
	"bytes"
	"strings"
)

// Truncate the given string to length using … as ellipsis.
func Truncate(s string, length int) string {
	return TruncateWithEllipsis(s, length, "…")
}

// TruncateWithEllipsis truncates the given string to length using provided ellipsis.
func TruncateWithEllipsis(s string, length int, ellipsis string) string {

	l := len(s)
	el := len(ellipsis)
	if l+el > length {
		s = string(s[0:length-el]) + ellipsis
	}
	return s
}

// ToPlural provides the plural version of an English word using some simple rules and a table of exceptions.
func ToPlural(text string) (plural string) {

	// We only deal with lowercase
	word := strings.ToLower(text)

	// Check translations first, and return a direct translation if there is one
	if translations[word] != "" {
		return translations[word]
	}

	// If we have no translation, just follow some basic rules - avoid new rules if possible
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "z") || strings.HasSuffix(word, "h") {
		plural = word + "es"
	} else if strings.HasSuffix(word, "y") {
		plural = strings.TrimRight(word, "y") + "ies"
	} else if strings.HasSuffix(word, "um") {
		plural = strings.TrimRight(word, "um") + "a"
	} else {
		plural = word + "s"
	}

	return plural
}

// Which irregulars are important or correct depends on your usage of English
// Some of those below are now considered old-fashioned and many more could be added
// As this is used for database models, it only needs a limited subset of all irregulars
// NB you should not attempt to reverse and singularize, but just use the singular provided
var translations = map[string]string{
	"hero":        "heroes",
	"supernova":   "supernovae",
	"day":         "days",
	"monkey":      "monkeys",
	"money":       "monies",
	"chassis":     "chassis",
	"sheep":       "sheep",
	"aircraft":    "aircraft",
	"fish":        "fish",
	"nucleus":     "nuclei",
	"mouse":       "mice",
	"buffalo":     "buffalo",
	"species":     "species",
	"information": "information",
	"wife":        "wives",
	"shelf":       "shelves",
	"index":       "indices",
	"matrix":      "matrices",
	"formula":     "formulae",
	"millennium":  "millennia",
	"ganglion":    "ganglia",
	"octopus":     "octopodes",
	"man":         "men",
	"woman":       "women",
	"person":      "people",
	"axis":        "axes",
	"die":         "dice",
	// ..etc
}

// ToSnake converts a string from struct field names to corresponding database column names (e.g. FieldName to field_name)
func ToSnake(text string) string {
	b := bytes.NewBufferString("")
	for i, c := range text {
		if i > 0 && c >= 'A' && c <= 'Z' {
			b.WriteRune('_')
		}
		b.WriteRune(c)
	}
	return strings.ToLower(b.String())
}

// ToCamel converts a string from database column names to corresponding struct field names (e.g. field_name to FieldName)
func ToCamel(text string, private ...bool) string {
	lowerCamel := false
	if private != nil {
		lowerCamel = private[0]
	}
	b := bytes.NewBufferString("")
	s := strings.Split(text, "_")
	for i, v := range s {
		if len(v) > 0 {
			s := v[:1]
			if i > 0 || lowerCamel == false {
				s = strings.ToUpper(s)
			}
			b.WriteString(s)
			b.WriteString(v[1:])
		}
	}
	return b.String()
}
