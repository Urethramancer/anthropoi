package anthropoi_test

import (
	"testing"
	"time"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
)

const pw = "secret passphrase"

func TestSetPassword(t *testing.T) {
	u := anthropoi.User{}
	for cost := 10; cost < 14; cost++ {
		start := time.Now()
		err := u.SetPassword(pw, cost)
		stop := time.Now()
		if err != nil {
			t.Logf("Error setting password: %s", err.Error())
			t.Fail()
		} else {
			t.Logf("Cost %d took %v to generate %s%s%s from %s%s%s\n",
				cost, stop.Sub(start), ansi.Green, u.Password, ansi.Normal, ansi.Green, pw, ansi.Normal)
		}
	}
}

func TestSetDovecotPassword(t *testing.T) {
	multi := 10000
	salt := anthropoi.GenString(16)
	for i := 5; i < 17; i++ {
		start := time.Now()
		hash := anthropoi.GenerateDovecotPassword(pw, salt, i*multi)
		stop := time.Now()
		t.Logf("%d iterations took %v to generate %s%s%s from %s%s%s\n",
			i*multi, stop.Sub(start), ansi.Green, hash, ansi.Normal, ansi.Green, pw, ansi.Normal)
	}
}
