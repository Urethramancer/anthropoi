package anthropoi_test

import (
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
)

const pw = "secret passphrase"

func TestAcceptablePassword(t *testing.T) {
	u := anthropoi.User{
		Username: "pierredelecto",
		First:    "Pierre",
		Last:     "Delecto",
		Email:    "pierre@delecto.xxx",
	}

	if u.AcceptablePassword("12345") {
		t.Logf("Validity check failed: Fully numeric passwords should not work!")
		t.FailNow()
	} else {
		t.Log("Numeric password not accepted. Good.")
	}

	if u.AcceptablePassword("pierre") {
		t.Logf("Validity check failed: Name as password should not work!")
		t.FailNow()
	} else {
		t.Log("Name & password similarity is considered invalid. Good.")
	}
}

func TestGetCost(t *testing.T) {
	pw, err := bcrypt.GenerateFromPassword([]byte("moo"), 10)
	if err != nil {
		t.Logf("Failed to generate bcrypt password.")
		t.FailNow()
	}
	u := &anthropoi.User{Password: string(pw)}
	c := u.GetCost()
	if c == 10 {
		t.Logf("Password %s%s%s has cost %d. Good.", ansi.Green, pw, ansi.Normal, u.GetCost())
	} else {
		t.Logf("Cost mismatch!")
		t.FailNow()
	}
}

func TestGetRounds(t *testing.T) {
	u := anthropoi.User{Password: anthropoi.GenerateDovecotPassword("moo", "salt", 10000)}
	r := u.GetRounds()
	if r == 10000 {
		t.Logf("Passwword %s%s%s has %d rounds. Good.", ansi.Green, u.Password, ansi.Normal, r)
	} else {
		t.Logf("Rounds mismatch.")
		t.FailNow()
	}
}

func TestSetPassword(t *testing.T) {
	u := anthropoi.User{}
	for cost := 10; cost < 14; cost++ {
		start := time.Now()
		err := u.SetPassword(pw, cost)
		stop := time.Now()
		if err != nil {
			t.Logf("Error setting password: %s", err.Error())
			t.FailNow()
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

func TestCompareDovecotPassword(t *testing.T) {
	salt := "0123456789abcdef"
	password := "{SHA512-CRYPT}$6$rounds=10000$0123456789abcdef$yWg2ncsjJEyAkbcwd.XkLNHpdZ30gK4QX9YWC1mds1pL7noAF.6Xly7VM1X8BLCCmZjt2IFGz8f8EiU44bjNf/"
	u := anthropoi.User{
		Salt:     salt,
		Password: password,
	}
	if !u.CompareDovecotHashAndPassword(pw) {
		t.Logf("Error comparing passwords!")
		t.FailNow()
	} else {
		t.Logf("Password matches salt and pre-generated hash.")
	}
}
