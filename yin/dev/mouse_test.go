package dev_test

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/madlambda/qi/yin/dev"
	"github.com/madlambda/spells/assert"
)

func parseStatus(t *testing.T, content []byte) (int, int) {
	parts := bytes.Split(content, []byte(" "))
	x, err := strconv.ParseInt(string(parts[0]), 10, 32)
	assert.NoError(t, err)

	y, err := strconv.ParseInt(string(parts[1]), 10, 32)
	assert.NoError(t, err)

	return int(x), int(y)
}

func TestMouseUpdate(t *testing.T) {
	for _, tc := range []struct {
		ex, ey int
	}{
		{0, 0},
		{10, 100},
		{30, 200},
	} {
		testMouseUpdate(t, tc.ex, tc.ey)
	}
}

func testMouseUpdate(t *testing.T, ex, ey int) {
	m := dev.MouseInit()
	done := make(chan bool)
	go func() {
		done <- true
		m.UpdateCoords(ex, ey)
	}()

	<-done
	content, err := ioutil.ReadAll(m)
	assert.NoError(t, err)

	x, y := parseStatus(t, content)
	if x != ex || y != ey {
		t.Errorf("expected (%d, %d) but got (%d, %d)",
			ex, ey, x, y)
	}
}
