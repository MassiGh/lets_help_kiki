package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ProblemSolver func(FirstLineInput, []PackageDetail, [][]string) ([]string, error)
type Problem struct {
	Key        string
	Title      string
	ExtraLines int
	Solver     ProblemSolver
}

func main() {
	// List of problems
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

	reader := bufio.NewReader(os.Stdin)

	// Get problem
	problem := pickProblem(reader, problems)
	// Get problems info
	firstLineInput, packageDetails, extraDetails := readProblemInputs(reader, problem)
	// Solve the problem
	outputs, err := problem.Solver(firstLineInput, packageDetails, extraDetails)
	if err != nil {
		fmt.Println(err)
	}
	// Write the outputs in console
	fmt.Println("<----------- Output ----------->")
	for _, output := range outputs {
		fmt.Println(output)
	}
}

// Function to show options to the user to select one of the problems
func pickProblem(reader *bufio.Reader, problems []Problem) Problem {
	displayProblems(problems)

	problem, err := getSelectedProblem(reader, problems)
	if err != nil {
		fmt.Println(err)
		pickProblem(reader, problems)
	}

	printData := fmt.Sprintf("<----------- Selected problem: %s ----------->", problem.Title)
	fmt.Println(printData)
	return problem
}

// Function to just show the list of problems in the console
func displayProblems(problems []Problem) {
	fmt.Println("<----------- What do you want me to calculate? ----------->")
	for _, problem := range problems {
		printData := fmt.Sprintf("%s. %s", problem.Key, problem.Title)
		fmt.Println(printData)
	}
	fmt.Println("<----------- Please enter your choice number and press Enter: ----------->")
}

// Function to read the problem number from stdin and return the selected problem
func getSelectedProblem(reader *bufio.Reader, problems []Problem) (Problem, error) {
	problemNumber := strings.TrimSpace(readLine(reader))

	for _, problem := range problems {
		if problemNumber == problem.Key {
			return problem, nil
		}
	}

	return Problem{}, fmt.Errorf("read problem error: '%s' is not a known problem number", problemNumber)
}

// Function to read the problem inputs from stdin, validate and parse them
// extra detail line is just read in this function and validatation is handled in the solver function
func readProblemInputs(reader *bufio.Reader, problem Problem) (FirstLineInput, []PackageDetail, [][]string) {
	fmt.Println("<----------- Please enter base cost and number of packages ----------->")
	firstLineInput := getFirstLineInput(reader)

	printData := fmt.Sprintf("<----------- Please enter %d package details and extra inputs ----------->", firstLineInput.NumberOfPackages)
	fmt.Println(printData)
	packageDetails := []PackageDetail{}
	for len(packageDetails) < firstLineInput.NumberOfPackages {
		packageDetail := getPackageDetail(reader, len(packageDetails))
		packageDetails = append(packageDetails, packageDetail)
	}

	extraDetails := [][]string{}
	if problem.ExtraLines > 0 {
		for len(extraDetails) < problem.ExtraLines {
			inputTokens := strings.Split(strings.TrimSpace(readLine(reader)), " ")
			extraDetails = append(extraDetails, inputTokens)
		}
	}

	return firstLineInput, packageDetails, extraDetails
}

// Function to read first line of input from stdin
func getFirstLineInput(reader *bufio.Reader) FirstLineInput {
	inputTokens := strings.Split(strings.TrimSpace(readLine(reader)), " ")
	firstLineInput, err := parseFirstLineInput(inputTokens)
	if err != nil {
		fmt.Println(err)
		getFirstLineInput(reader)
	}
	return firstLineInput
}

// Function to validate and parse the first line of inputs
// first line of input should have two integer "baseCost(int) numberOfPackages(int)"
func parseFirstLineInput(inputTokens []string) (FirstLineInput, error) {
	var firstLineInput FirstLineInput
	if len(inputTokens) != 2 {
		return firstLineInput, fmt.Errorf("parse first input line error: Wrong number of inputs")
	}

	baseCost, err := strconv.Atoi(inputTokens[0])
	if err != nil {
		return firstLineInput, fmt.Errorf("parse first input line error: Wrong base cost input")
	}

	numberOfPackages, err := strconv.Atoi(inputTokens[1])
	if err != nil {
		return firstLineInput, fmt.Errorf("parse first input line error: Wrong number of packages input")
	}

	firstLineInput = FirstLineInput{
		BaseCost:         baseCost,
		NumberOfPackages: numberOfPackages,
	}
	return firstLineInput, nil
}

// Function to read package details from stdin
func getPackageDetail(reader *bufio.Reader, index int) PackageDetail {
	inputTokens := strings.Split(strings.TrimSpace(readLine(reader)), " ")
	packageDetail, err := parsePackageDetail(inputTokens, index)
	if err != nil {
		fmt.Println(err)
		getPackageDetail(reader, index)
	}
	return packageDetail
}

// Function to parse package detail input "packageId(string) weight(int) distance(int) offerIds(comma seperated string)"
func parsePackageDetail(inputTokens []string, index int) (PackageDetail, error) {
	var packageDetail PackageDetail
	if len(inputTokens) != 4 {
		return packageDetail, fmt.Errorf("parse package inputs error: Wrong number of inputs")
	}

	weight, err := strconv.Atoi(inputTokens[1])
	if err != nil {
		return packageDetail, fmt.Errorf("parse package inputs error: Wrong package weight input")
	}

	distance, err := strconv.Atoi(inputTokens[2])
	if err != nil {
		return packageDetail, fmt.Errorf("parse package inputs error: Wrong package distance input")
	}

	offerIds := strings.Split(inputTokens[3], ",")

	packageDetail = PackageDetail{
		Index:    index,
		Title:    inputTokens[0],
		Weight:   weight,
		Distance: distance,
		OfferIds: offerIds,
	}
	return packageDetail, nil
}

// Function to read a line from stdin
func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}
