package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transliterate(t *testing.T) {
	assert.EqualValues(t, transliterate("hello"), "hello")
	assert.EqualValues(t, transliterate("שלום"), "shlom")
	assert.EqualValues(t, transliterate("שלום רב שובך ציפורה נחמדת"), "shlom rb shobkh tsyporh n7mdt")
}
