package main

import "fmt"

type CalculationOutput struct {
	TotalCost int
	Discount  int
}

type CompareAmount struct {
	GreaterThanEqual int
	LessThanEqual    int
}

type Offer struct {
	Id       string
	Distance CompareAmount
	Weight   CompareAmount
	Percent  int
}

// The solver function for the "Delivery Cost Estimation" problem
func CalculateDeliveryCost(firstInputLine FirstLineInput, packageDetails []PackageDetail, extraDetails [][]string) ([]string, error) {
	baseDeliveryCost := firstInputLine.BaseCost

	outputs := []string{}
	for _, packageDetail := range packageDetails {
		calculationOutput := calculateTotalCost(baseDeliveryCost, packageDetail)
		outputs = append(outputs, fmt.Sprintf("%s %d %d", packageDetail.Title, calculationOutput.Discount, calculationOutput.TotalCost))
	}
	return outputs, nil
}

// Function to calculate the total cost according to the base cost, weight, distance and discounts
func calculateTotalCost(baseDeliveryCost int, packageDetail PackageDetail) CalculationOutput {
	deliveryCost := baseDeliveryCost + (packageDetail.Weight * 10) + (packageDetail.Distance * 5)
	discount := calculateDiscounts(packageDetail, deliveryCost)

	return CalculationOutput{
		TotalCost: deliveryCost - discount,
		Discount:  discount,
	}
}

// Function to calculate discount for the package
func calculateDiscounts(packageDetail PackageDetail, deliveryCost int) int {
	offers := []Offer{
		{
			Id: "OFR001",
			Distance: CompareAmount{
				GreaterThanEqual: 0,
				LessThanEqual:    199,
			},
			Weight: CompareAmount{
				GreaterThanEqual: 70,
				LessThanEqual:    200,
			},
			Percent: 10,
		},
		{
			Id: "OFR002",
			Distance: CompareAmount{
				GreaterThanEqual: 50,
				LessThanEqual:    150,
			},
			Weight: CompareAmount{
				GreaterThanEqual: 100,
				LessThanEqual:    250,
			},
			Percent: 7,
		},
		{
			Id: "OFR003",
			Distance: CompareAmount{
				GreaterThanEqual: 50,
				LessThanEqual:    250,
			},
			Weight: CompareAmount{
				GreaterThanEqual: 10,
				LessThanEqual:    150,
			},
			Percent: 5,
		},
	}

	discount := 0
	for _, offerId := range packageDetail.OfferIds {
		for _, offer := range offers {
			if offer.Id == offerId {
				if (offer.Weight.GreaterThanEqual == 0 || packageDetail.Weight >= offer.Weight.GreaterThanEqual) &&
					(offer.Weight.LessThanEqual == 0 || packageDetail.Weight <= offer.Weight.LessThanEqual) &&
					(offer.Distance.GreaterThanEqual == 0 || packageDetail.Distance >= offer.Distance.GreaterThanEqual) &&
					(offer.Distance.LessThanEqual == 0 || packageDetail.Distance <= offer.Distance.LessThanEqual) {

					discount += deliveryCost * offer.Percent / 100
				}
			}
		}
	}
	return discount
}
