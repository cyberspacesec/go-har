package har

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHar(t *testing.T) {
	harFileContent, err := os.ReadFile("testdata/example.har")
	assert.Nil(t, err)
	har, err := ParseHar(harFileContent)
	assert.Nil(t, err)
	t.Log(har)
}

func TestParseHarFile(t *testing.T) {
	har, err := ParseHarFile("testdata/example.har")
	assert.Nil(t, err)
	t.Log(har)
}
