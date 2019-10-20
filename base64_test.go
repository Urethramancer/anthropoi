package anthropoi_test

import (
	"strings"
	"testing"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
)

func TestBase6424(t *testing.T) {
	in := "this is a test string"
	desired := "oVKOn/GOn/GMUELNnF56nFbQdtqN"
	out := string(anthropoi.Base6424(in))

	if strings.Compare(out, desired) != 0 {
		t.Logf("Output mismatch:\n\t'%s%s%s' (%d) is not\n\t'%s%s%s' (%d)", ansi.Green, out, ansi.Normal, len(out), ansi.Green, desired, ansi.Normal, len(desired))
		t.Fail()
	} else {
		t.Logf("'%s' = '%s'", out, desired)
	}
}
