package main

type FirstLineInput struct {
	BaseCost         int
	NumberOfPackages int
}

type PackageDetail struct {
	Index    int
	Title    string
	Weight   int
	Distance int
	OfferIds []string
}
