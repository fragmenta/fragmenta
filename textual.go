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
	"addendum":    "addenda", // or addendums
	"aircraft":    "aircraft",
	"alumna":      "alumnae",
	"alumnus":     "alumni",
	"analysis":    "analyses",
	"antenna":     "antennae", // or antennas
	"antithesis":  "antitheses",
	"apex":        "apices",     // or apexes
	"appendix":    "appendices", // or appendixes
	"axis":        "axes",
	"bacillus":    "bacilli",
	"bacterium":   "bacteria",
	"basis":       "bases",
	"beau":        "beaux",   // or beaus
	"belief":      "beliefs", // not irregular?
	"bison":       "bison",
	"buffalo":     "buffalo",
	"bureau":      "bureaux", // or bureaus
	"cactus":      "cacti",   // or cactus or cactuses
	"chassis":     "chassis",
	"chateau":     "chateaux", // or chateaus
	"château":     "châteaux", // or châteaus
	"child":       "children",
	"codex":       "codices",
	"concerto":    "concerti", // or concertos
	"corpus":      "corpora",
	"crisis":      "crises",
	"criterion":   "criteria",  // or criterions
	"curriculum":  "curricula", // or curriculums
	"datum":       "data",
	"day":         "days",
	"deer":        "deer", // or deers
	"die":         "dice",
	"diagnosis":   "diagnoses",
	"die":         "dice",    // or dies
	"dwarf":       "dwarves", // or dwarfs
	"ellipsis":    "ellipses",
	"erratum":     "errata",
	"faux pas":    "faux pas",
	"fez":         "fezzes",   // or fezes
	"fish":        "fish",     // or fishes
	"focus":       "foci",     // or focuses
	"foot":        "feet",     // or foot
	"formula":     "formulae", // or formulas
	"fungus":      "fungi",    // or funguses
	"ganglion":    "ganglia",
	"genus":       "genera", // or genuses
	"goose":       "geese",
	"graffito":    "graffiti",
	"grouse":      "grouse", // or grouses
	"half":        "halves",
	"hero":        "heroes",
	"hoof":        "hooves", // or hoofs
	"hypothesis":  "hypotheses",
	"index":       "indices", // or indexes
	"information": "information",
	"larva":       "larvae",   // or larvas
	"libretto":    "libretti", // or librettos
	"life":        "lives",
	"loaf":        "loaves",
	"locus":       "loci",
	"louse":       "lice",
	"man":         "men",
	"matrix":      "matrices",  // or matrixes
	"medium":      "media",     // or mediums
	"memorandum":  "memoranda", // or memorandums
	"millennium":  "millennia",
	"minutia":     "minutiae",
	"moose":       "moose",
	"money":       "monies", // not irregular?
	"monkey":      "monkeys",
	"mouse":       "mice",
	"nebula":      "nebulae", // or nebulas
	"nucleus":     "nuclei",  // or nucleuses
	"oasis":       "oases",
	"offspring":   "offspring", // or offsprings
	"opus":        "opera",     // or opuses
	"octopus":     "octopodes",
	"ovum":        "ova",
	"ox":          "oxen", // or ox
	"parenthesis": "parentheses",
	"person":      "people",
	"phenomenon":  "phenomena", // or phenomenons
	"phylum":      "phyla",
	"prognosis":   "prognoses",
	"quiz":        "quizzes",
	"radius":      "radii",     // or radiuses
	"referendum":  "referenda", // or referendums
	"salmon":      "salmon",    // or salmons
	"scarf":       "scarves",   // or scarfs
	"self":        "selves",
	"series":      "series",
	"sheep":       "sheep",
	"shrimp":      "shrimp", // or shrimps
	"shelf":       "shelves",
	"species":     "species",
	"stimulus":    "stimuli",
	"stratum":     "strata",
	"swine":       "swine",
	"syllabus":    "syllabi",  // or syllabuses
	"symposium":   "symposia", // or symposiums
	"synopsis":    "synopses",
	"supernova":   "supernovae",
	"tableau":     "tableaux", // or tableaus
	"thesis":      "theses",
	"thief":       "thieves",
	"tooth":       "teeth",
	"trout":       "trout",     // or trouts
	"tuna":        "tuna",      // or tunas
	"vertebra":    "vertebrae", // or vertebras
	"vertex":      "verticies", // or vertexes
	"vita":        "vitae",
	"vortex":      "vorticies", // or vortexes
	"wharf":       "wharves",   // or wharfs
	"wife":        "wives",
	"wolf":        "wolves",
	"woman":       "women",
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
