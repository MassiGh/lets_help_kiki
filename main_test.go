package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPickProblem(t *testing.T) {
	problems := []Problem{
		{
			Key:        "1",
			Title:      "Delivery Cost Estimation with Offers",
			ExtraLines: 0,
			Solver:     CalculateDeliveryCost,
		},
		{
			Key:        "2",
			Title:      "Delivery Time Estimation",
			ExtraLines: 1,
			Solver:     CalculateDeliveryTime,
		},
	}
	t.Run("return problem for the valid problem number", func(t *testing.T) {
		problemNumber := "1"
		reader := bufio.NewReader(strings.NewReader(problemNumber))
		problem := pickProblem(reader, problems)

		assert.Equal(t, "Delivery Cost Estimation with Offers", problem.Title)
	})
}

func TestGetSelectedProblem(t *testing.T) {
	problems := []Problem{
		{
			Key:        "1",
			Title:      "Delivery Cost Estimation with Offers",
			ExtraLines: 0,
			Solver:     CalculateDeliveryCost,
		},
		{
			Key:        "2",
			Title:      "Delivery Time Estimation",
			ExtraLines: 1,
			Solver:     CalculateDeliveryTime,
		},
	}
	t.Run("return error for empty problem number", func(t *testing.T) {
		problemNumber := ""
		reader := bufio.NewReader(strings.NewReader(problemNumber))
		problem, err := getSelectedProblem(reader, problems)

		assert.Error(t, err)
		assert.Equal(t, Problem{}, problem)
	})
	t.Run("return error for the invalid problem number", func(t *testing.T) {
		problemNumber := "186"
		reader := bufio.NewReader(strings.NewReader(problemNumber))
		problem, err := getSelectedProblem(reader, problems)

		assert.Error(t, err)
		assert.Equal(t, Problem{}, problem)
	})
	t.Run("return problem for the valid problem number", func(t *testing.T) {
		problemNumber := "1"
		reader := bufio.NewReader(strings.NewReader(problemNumber))
		problem, err := getSelectedProblem(reader, problems)

		assert.NoError(t, err)
		assert.Equal(t, "Delivery Cost Estimation with Offers", problem.Title)
	})
}

func TestGetFirstLineInput(t *testing.T) {
	t.Run("return firstInputLine for the valid input", func(t *testing.T) {
		firstLineInput := "100 5"
		reader := bufio.NewReader(strings.NewReader(firstLineInput))
		inputTokens := getFirstLineInput(reader)

		assert.Equal(t, FirstLineInput{
			BaseCost:         100,
			NumberOfPackages: 5,
		}, inputTokens)
	})
}

func TestParseFirstInputLine(t *testing.T) {
	t.Run("doesn't check wrong number of inputs in the first input line", func(t *testing.T) {
		inputTokens := []string{"100", "3", "2"}
		firstLineInput, err := parseFirstLineInput(inputTokens)

		assert.Error(t, err)
		assert.Equal(t, FirstLineInput{}, firstLineInput)
	})
	t.Run("doesn't check wrong base cost in the first input line", func(t *testing.T) {
		inputTokens := []string{"10s", "3"}
		firstLineInput, err := parseFirstLineInput(inputTokens)

		assert.Error(t, err)
		assert.Equal(t, FirstLineInput{}, firstLineInput)
	})
	t.Run("doesn't check wrong number of packages in the first input line", func(t *testing.T) {
		inputTokens := []string{"100", "ddd"}
		firstLineInput, err := parseFirstLineInput(inputTokens)

		assert.Error(t, err)
		assert.Equal(t, FirstLineInput{}, firstLineInput)
	})
	t.Run("return the first input line object", func(t *testing.T) {
		inputTokens := []string{"100", "3"}
		firstLineInput, err := parseFirstLineInput(inputTokens)

		assert.NoError(t, err)
		assert.Equal(t, FirstLineInput{BaseCost: 100, NumberOfPackages: 3}, firstLineInput)
	})
}

func TestGetPackageDetails(t *testing.T) {
	t.Run("return packageDetail for the valid input", func(t *testing.T) {
		packageDetailsInput := "PKG1 50 30 OFR001"
		reader := bufio.NewReader(strings.NewReader(packageDetailsInput))
		inputTokens := getPackageDetail(reader, 0)

		assert.Equal(t, PackageDetail{
			Index:    0,
			Title:    "PKG1",
			Weight:   50,
			Distance: 30,
			OfferIds: []string{"OFR001"},
		}, inputTokens)
	})
}

func TestParsePackageDetail(t *testing.T) {
	t.Run("doesn't check wrong number of inputs in the package detail input", func(t *testing.T) {
		inputTokens := []string{"PKG1", "50", "30"}
		packageDetail, err := parsePackageDetail(inputTokens, 0)

		assert.Error(t, err)
		assert.Equal(t, PackageDetail{}, packageDetail)
	})
	t.Run("doesn't check wrong weight in the package detail input", func(t *testing.T) {
		inputTokens := []string{"PKG1", "50s", "30", "OFR001"}
		packageDetail, err := parsePackageDetail(inputTokens, 0)

		assert.Error(t, err)
		assert.Equal(t, PackageDetail{}, packageDetail)
	})
	t.Run("doesn't check wrong distance in the package detail input", func(t *testing.T) {
		inputTokens := []string{"PKG1", "50", "30s", "OFR001"}
		packageDetail, err := parsePackageDetail(inputTokens, 0)

		assert.Error(t, err)
		assert.Equal(t, PackageDetail{}, packageDetail)
	})
	t.Run("return the package detail object", func(t *testing.T) {
		inputTokens := []string{"PKG1", "50", "30", "OFR001"}
		packageDetail, err := parsePackageDetail(inputTokens, 0)

		assert.NoError(t, err)
		assert.Equal(t, PackageDetail{
			Index:    0,
			Title:    "PKG1",
			Weight:   50,
			Distance: 30,
			OfferIds: []string{"OFR001"},
		}, packageDetail)
	})
}
