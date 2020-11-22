package main

import (
	"fmt"
	"strings"
	"testing"
)

var results string = `ok  	git.x.com/srv/x/src/reviews/actions	(cached)
ok  	git.x.com/srv/x/src/search/actions	(cached)
?  	git.x.com/srv/x/src/seasons/actions [no test files]
--- FAIL: TestName (2.25s)
    file_test.go:174: test: error in template: x/views/x.html.got:3: function "x" not defined
FAIL
FAIL	git.x.com/srv/x/src/sites/actions	2.619s
FAIL	git.x.com/srv/x/src/bookword	2.619s
ok  	git.x.com/srv/x/src/tags/actions	(cached)
FAIL`

func TestColorizeResults(t *testing.T) {
	output := colorizeResults(results)

	// Check ok is green
	if !strings.Contains(output, ColorGreen+"ok") {
		t.Errorf("colorise: ok failed:\n%s", output)
	}

	// Check ok within a word is not green
	if strings.Contains(output, ColorGreen+"okword") {
		t.Errorf("colorise: ok within word incorrectly green:\n%s", output)
	}

	// Check empty fail removed
	if strings.HasSuffix(output, "FAIL") {
		t.Errorf("colorise: empty fail not removed:\n%s", output)
	}

	// Check ? is amber
	if !strings.Contains(output, ColorAmber+"?") {
		t.Errorf("colorise: ? failed:\n%s", output)
	}

	// Check FAIL is red
	if !strings.Contains(output, ColorRed+"FAIL") {
		t.Errorf("colorise: FAIL failed:\n%s", output)
	}

	// Check 	--- FAIL: is red
	if !strings.Contains(output, ColorRed+"--- FAIL:") {
		t.Errorf("colorise: FAIL failed:\n%s", output)
	}

	fmt.Println(results)
}
