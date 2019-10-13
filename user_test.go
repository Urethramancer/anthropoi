package anthropoi_test

import (
	"testing"
	"time"

	"github.com/Urethramancer/anthropoi"
)

func TestSetPassword(t *testing.T) {
	u := anthropoi.User{}
	for cost := 10; cost < 16; cost++ {
		start := time.Now()
		err := u.SetPassword("secret passphrase", cost)
		if err != nil {
			t.Logf("Error setting password: %s", err.Error())
			t.Fail()
		} else {
			stop := time.Now()
			t.Logf("Cost %d took %v\n", cost, stop.Sub(start))
		}
	}
}
