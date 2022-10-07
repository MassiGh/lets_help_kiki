package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDeliveryTime(t *testing.T) {
	firstLineInput := FirstLineInput{
		BaseCost:         100,
		NumberOfPackages: 5,
	}

	extraDetails := [][]string{{"2", "70", "200"}}

	t.Run("return error if extraDetails are empty", func(t *testing.T) {
		outputs, err := CalculateDeliveryTime(firstLineInput, []PackageDetail{}, nil)

		assert.Error(t, err, "Validate extra details error: Wrong number of inputs")
		assert.Equal(t, []string(nil), outputs)
	})

	t.Run("return correct output", func(t *testing.T) {
		packageDetails := []PackageDetail{
			{
				Index:    0,
				Title:    "PKG1",
				Weight:   50,
				Distance: 30,
				OfferIds: []string{"OFR001"},
			},
			{
				Index:    1,
				Title:    "PKG2",
				Weight:   75,
				Distance: 125,
				OfferIds: []string{"OFR008"},
			},
			{
				Index:    2,
				Title:    "PKG3",
				Weight:   175,
				Distance: 100,
				OfferIds: []string{"OFR003"},
			},
			{
				Index:    3,
				Title:    "PKG4",
				Weight:   110,
				Distance: 60,
				OfferIds: []string{"OFR002"},
			},
			{
				Index:    4,
				Title:    "PKG5",
				Weight:   155,
				Distance: 95,
				OfferIds: []string{"NA"},
			},
		}
		outputs, err := CalculateDeliveryTime(firstLineInput, packageDetails, extraDetails)

		assert.NoError(t, err)
		assert.Equal(t, []string{"PKG1 0 750 3.98", "PKG2 0 1475 1.78", "PKG3 0 2350 1.42", "PKG4 105 1395 0.85", "PKG5 0 2125 4.19"}, outputs)
	})
}
