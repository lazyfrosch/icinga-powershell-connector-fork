package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseNetstring(t *testing.T) {
	r, err := os.Open("testdata/netstring.dat")
	if err != nil {
		t.Fatal(err)
	}

	data, err := ParseNetstring(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte("a"), data)

	data, err = ParseNetstring(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte("bbbbb"), data)

	data, err = ParseNetstring(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte("aaaaaaaaaaa"), data)

	data, err = ParseNetstring(r)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
