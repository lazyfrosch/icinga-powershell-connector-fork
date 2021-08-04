package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIPerfdataList_UnmarshalJSON(t *testing.T) {
	var pl APIPerfdataList

	err := json.Unmarshal([]byte("{}"), &pl)
	assert.NoError(t, err)
	assert.Equal(t, APIPerfdataList{}, pl)

	err = json.Unmarshal([]byte(`["a", "b", "c"]`), &pl)
	assert.NoError(t, err)
	assert.Equal(t, APIPerfdataList{"a", "b", "c"}, pl)
}
