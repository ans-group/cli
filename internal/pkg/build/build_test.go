package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfo_String_ExpectedOutput(t *testing.T) {
	t.Run("PopulatedProperties", func(t *testing.T) {
		b := BuildInfo{
			Version:   "v1.0.0",
			BuildDate: "02-01-2019",
		}

		out := b.String()

		assert.Equal(t, "v1.0.0 built on 02-01-2019", out)
	})
	t.Run("MissingProperties", func(t *testing.T) {
		b := BuildInfo{}

		out := b.String()

		assert.Equal(t, "UNKNOWN built on UNKNOWN", out)
	})
}
