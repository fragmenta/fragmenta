package main

import (
	"bytes"
	"strings"
)

// ToPlural provides the plural version of an English word
// using some simple rules and a table of exceptions.
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

// ToCamel converts a string from database column names
// to corresponding struct field names
// e.g. field_name to FieldName.
func ToCamel(s string) string {

	b := bytes.NewBufferString("")
	words := strings.Split(s, "_")
	for _, word := range words {

		// If the word matches commonInitialisms like ID or HTML write as uppercase
		if commonInitialisms[word] {
			b.WriteString(strings.ToUpper(word))
			continue
		}

		// Ignore zero length words
		if len(word) == 0 {
			continue
		}

		// Write the initial cap
		b.WriteString(strings.ToUpper(word[:1]))

		// Write the rest of the word
		b.WriteString(word[1:])

	}
	return b.String()
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

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
// this set of initialisms is taken from golang.org/tools/lint
// to avoid generating field names which the linter will complain about.
var commonInitialisms = map[string]bool{
	"acl":   true,
	"api":   true,
	"ascii": true,
	"cpu":   true,
	"css":   true,
	"dns":   true,
	"eof":   true,
	"guid":  true,
	"html":  true,
	"http":  true,
	"https": true,
	"id":    true,
	"ip":    true,
	"json":  true,
	"lhs":   true,
	"qps":   true,
	"ram":   true,
	"rhs":   true,
	"rpc":   true,
	"sla":   true,
	"smtp":  true,
	"sql":   true,
	"ssh":   true,
	"tcp":   true,
	"tls":   true,
	"ttl":   true,
	"udp":   true,
	"ui":    true,
	"uid":   true,
	"uuid":  true,
	"uri":   true,
	"url":   true,
	"utf8":  true,
	"vm":    true,
	"xml":   true,
	"xmpp":  true,
	"xsrf":  true,
	"xss":   true,
}
