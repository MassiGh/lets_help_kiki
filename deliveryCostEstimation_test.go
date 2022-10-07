package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDeliveryCost(t *testing.T) {
	firstLineInput := FirstLineInput{
		BaseCost:         100,
		NumberOfPackages: 3,
	}

	t.Run("return zero for invalid offerId", func(t *testing.T) {
		packageDetails := []PackageDetail{
			{
				Index:    0,
				Title:    "PKG1",
				Weight:   5,
				Distance: 5,
				OfferIds: []string{"OFR001"},
			},
		}

		outputs, err := CalculateDeliveryCost(firstLineInput, packageDetails, nil)

		assert.NoError(t, err)
		assert.Equal(t, []string{"PKG1 0 175"}, outputs)
	})

	t.Run("return correct discount for valid offerId", func(t *testing.T) {
		packageDetails := []PackageDetail{
			{
				Index:    0,
				Title:    "PKG3",
				Weight:   10,
				Distance: 100,
				OfferIds: []string{"OFR003"},
			},
		}
		outputs, err := CalculateDeliveryCost(firstLineInput, packageDetails, nil)

		assert.NoError(t, err)
		assert.Equal(t, []string{"PKG3 35 665"}, outputs)
	})
}
