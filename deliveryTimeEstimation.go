package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

type Subset struct {
	TotalWeight          int   // The sum of packages weight in the subset
	MaxDistance          int   // The maximum distance of the packages in the subset
	PackageDetailIndices []int // List of ids for removing them from the main array after shipment
	PackageDetails       []PackageDetail
}

type ExtraDetails struct {
	NumberOfVehicles   int
	MaxSpeed           int
	MaxCarriableWeight int
}

type ShipmentDetail struct {
	Title        string
	Discount     int
	TotalCost    int
	DeliveryTime float64
}

// The solver function for the "Delivery Time Estimation" problem
func CalculateDeliveryTime(firstInputLine FirstLineInput, packageDetails []PackageDetail, extraDetails [][]string) ([]string, error) {

	validatedExtraDetails, err := validateExtraDetails(extraDetails)
	if err != nil {
		return nil, err
	}

	sort.Slice(packageDetails, func(i, j int) bool {
		return packageDetails[i].Weight < packageDetails[j].Weight
	})

	shipmentSubsets := make([]Subset, 0)
	getShipmentSubsets(packageDetails, validatedExtraDetails.MaxCarriableWeight, &shipmentSubsets)

	shipmentDetails := calculateShipmentDetails(firstInputLine, shipmentSubsets, validatedExtraDetails)
	outputs := []string{}
	for _, o := range shipmentDetails {
		outputs = append(outputs, fmt.Sprintf("%s %d %d %.2f", o.Title, o.Discount, o.TotalCost, o.DeliveryTime))
	}

	return outputs, nil
}

// Function to get all shipment subsets
func getShipmentSubsets(packages []PackageDetail, maxCarriableWeight int, result *[]Subset) {
	maxSize := findMaxSubsetSize(packages, maxCarriableWeight)

	current := Subset{}
	subsets := make([]Subset, 0)

	getSubsets(packages, maxSize, 0, current, maxCarriableWeight, &subsets)

	bestSubset := getBestSubset(subsets)

	remainedPackages := removePackage(packages, bestSubset.PackageDetailIndices)
	*result = append(*result, bestSubset)
	if len(remainedPackages) > 0 {
		getShipmentSubsets(remainedPackages, maxCarriableWeight, result)
	}
}

// Function to find the max possible subset size
// The given packages array should be sorted based on the weight
func findMaxSubsetSize(sortedPackages []PackageDetail, maxCarriableWeight int) int {
	sum := 0
	subsetSize := 0
	for i := 0; i < len(sortedPackages); i++ {
		sum += sortedPackages[i].Weight
		if sum <= maxCarriableWeight {
			subsetSize++
		} else {
			break
		}
	}
	return subsetSize
}

// Function to get all subsets with specific size of packages
// Note: the sum of packages weight in each subset should not exceed maxCarriableWeight
func getSubsets(packages []PackageDetail, subsetSize int, index int, currentSubset Subset, maxCarriableWeight int, subsets *[]Subset) {
	// successful stop clause
	if len(currentSubset.PackageDetailIndices) == subsetSize {
		*subsets = append(*subsets, currentSubset)
		return
	}

	// unseccessful stop clause
	if index == len(packages) {
		return
	}

	currentPackage := packages[index]
	nextSubnet := currentSubset

	if currentSubset.TotalWeight+currentPackage.Weight <= maxCarriableWeight {
		currentSubset.PackageDetailIndices = append(currentSubset.PackageDetailIndices, index)
		currentSubset.PackageDetails = append(currentSubset.PackageDetails, packages[index])
		currentSubset.TotalWeight += currentPackage.Weight
		if currentSubset.MaxDistance < currentPackage.Distance {
			currentSubset.MaxDistance = currentPackage.Distance
		}
		// subset with x
		getSubsets(packages, subsetSize, index+1, currentSubset, maxCarriableWeight, subsets)
	}
	// subset without x
	getSubsets(packages, subsetSize, index+1, nextSubnet, maxCarriableWeight, subsets)
}

// Function to get the subset with max weight/distance between found subsets
func getBestSubset(subsets []Subset) Subset {
	if len(subsets) > 1 {
		sort.Slice(subsets, func(i, j int) bool {
			return (subsets[i].TotalWeight > subsets[j].TotalWeight ||
				(subsets[i].TotalWeight == subsets[j].TotalWeight && subsets[i].MaxDistance < subsets[j].MaxDistance))
		})
	}
	return subsets[0]
}

// Function to remove packages from array by index and return a new array
func removePackage(packages []PackageDetail, toRemoveIndexes []int) []PackageDetail {
	ret := make([]PackageDetail, 0)

	for i := 0; i < len(packages); i++ {
		isRemoved := false
		for _, j := range toRemoveIndexes {
			if i == j {
				isRemoved = true
				break
			}
		}
		if !isRemoved {
			ret = append(ret, packages[i])
		}
	}
	return ret
}

// Function to validate extra details based on the problem explenation
func validateExtraDetails(extraDetails [][]string) (ExtraDetails, error) {
	var extraDetail ExtraDetails
	if len(extraDetails) != 1 || len(extraDetails[0]) != 3 {
		return extraDetail, fmt.Errorf("%s", "Validate extra details error: Wrong number of inputs")
	}

	numberOfVehicles, err := strconv.Atoi(extraDetails[0][0])
	if err != nil {
		return extraDetail, fmt.Errorf("%s", "Validate extra details error: Wrong number of vehicles")
	}

	maxSpeed, err := strconv.Atoi(extraDetails[0][1])
	if err != nil {
		return extraDetail, fmt.Errorf("%s", "Validate extra details error: Wrong max speed")
	}

	maxCarriableWeight, err := strconv.Atoi(extraDetails[0][2])
	if err != nil {
		return extraDetail, fmt.Errorf("%s", "Validate extra details error: Wrong max carriable weight")
	}

	return ExtraDetails{
		NumberOfVehicles:   numberOfVehicles,
		MaxSpeed:           maxSpeed,
		MaxCarriableWeight: maxCarriableWeight,
	}, nil
}

// Function to calculate delivery time and create a shipmentDetail object for each package
func calculateShipmentDetails(firstInputLine FirstLineInput, shipmentSubsets []Subset, extraDetails ExtraDetails) []ShipmentDetail {

	result := make([]ShipmentDetail, firstInputLine.NumberOfPackages)
	vehiclesMaxDeliveryTime := make([]float64, extraDetails.NumberOfVehicles)

	for i := 0; i < len(shipmentSubsets); i++ {
		sort.Slice(vehiclesMaxDeliveryTime, func(i, j int) bool {
			return vehiclesMaxDeliveryTime[i] < vehiclesMaxDeliveryTime[j]
		})

		waitingTime := vehiclesMaxDeliveryTime[0]
		for _, d := range shipmentSubsets[i].PackageDetails {
			// Using the first problem
			deliveryCost := calculateTotalCost(firstInputLine.BaseCost, d)

			baseTime := float64(d.Distance) / float64(extraDetails.MaxSpeed)
			deliveryTime := roundoff(baseTime, 2) + waitingTime

			result[d.Index] = ShipmentDetail{
				Title:        d.Title,
				Discount:     deliveryCost.Discount,
				TotalCost:    deliveryCost.TotalCost,
				DeliveryTime: deliveryTime,
			}
		}
		maxDeliveryTime := float64(shipmentSubsets[i].MaxDistance) / float64(extraDetails.MaxSpeed)
		vehiclesMaxDeliveryTime[0] += (roundoff(maxDeliveryTime, 2) * 2)
	}
	return result
}

// Function to roundoff decimal points without rounding or flooring
func roundoff(num float64, floating_point float64) float64 {
	d := math.Pow(10, floating_point)
	return math.Floor(num*d) / d
}
