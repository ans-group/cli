package ecloud

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ecloudDescribeDescriptors(t *testing.T) {
	t.Run("NoPrefixCollisions", func(t *testing.T) {
		counts := make(map[string]int)
		for _, descriptor := range ecloudDescribeDescriptors {
			counts[descriptor.Prefix]++
		}

		for prefix, count := range counts {
			assert.Equal(t, 1, count, "duplicate descriptor for prefix [%s]", prefix)
		}
	})

	t.Run("AllDescriptorsHaveFetch", func(t *testing.T) {
		for _, descriptor := range ecloudDescribeDescriptors {
			assert.NotEmpty(t, descriptor.Prefix)
			assert.NotEmpty(t, descriptor.Name, "missing name for prefix [%s]", descriptor.Prefix)
			assert.NotNil(t, descriptor.Fetch, "missing fetch for prefix [%s]", descriptor.Prefix)

			for _, child := range descriptor.Children {
				assert.NotEmpty(t, child.Title, "missing child title for prefix [%s]", descriptor.Prefix)
				assert.NotNil(t, child.Fetch, "missing child fetch for prefix [%s]", descriptor.Prefix)
			}
		}
	})

	t.Run("PrefixesContainNoHyphen", func(t *testing.T) {
		for _, descriptor := range ecloudDescribeDescriptors {
			assert.False(t, strings.Contains(descriptor.Prefix, "-"), "prefix [%s] contains a hyphen", descriptor.Prefix)
		}
	})

	t.Run("AllDescriptorsIndexed", func(t *testing.T) {
		assert.Len(t, ecloudDescribeDescriptorsByPrefix, len(ecloudDescribeDescriptors))
	})
}

func Test_ecloudDescribePrefixes(t *testing.T) {
	t.Run("ReturnsSortedPrefixes", func(t *testing.T) {
		prefixes := ecloudDescribePrefixes()

		assert.Len(t, prefixes, len(ecloudDescribeDescriptors))
		assert.True(t, slices.IsSorted(prefixes))
		assert.Contains(t, prefixes, "i")
		assert.Contains(t, prefixes, "vpc")
	})
}
