package strategy

import (
	"errors"
)

type OrderPackageStrategy interface {
	ChoosePackage(price, weight int, additionalStretch bool) (int, int, error)
}

type BagPackageStrategy struct{}

func (s BagPackageStrategy) ChoosePackage(price, weight int, additionalStretch bool) (int, int, error) {
	if weight < 1 || price < 1 {
		return 0, 0, errors.New("price or weight can't be negetive or zero")
	}
	if weight > 10000 {
		return 0, weight, errors.New("to heavy for bag")
	}
	if additionalStretch {
		price += 1
	}
	return price + 5, weight, nil
}

type BoxPackageStrategy struct{}

func (s BoxPackageStrategy) ChoosePackage(price, weight int, additionalStretch bool) (int, int, error) {
	if weight < 1 || price < 1 {
		return 0, 0, errors.New("price or weight can't be negetive or zero")
	}
	if weight > 30000 {
		return 0, 0, errors.New("to heavy for box")
	}
	if additionalStretch {
		price += 1
	}
	return price + 20, weight, nil
}

type StretchPackageStrategy struct{}

func (s StretchPackageStrategy) ChoosePackage(price, weight int, additionalStretch bool) (int, int, error) {
	if additionalStretch {
		return 0, 0, errors.New("can't add stretch to stretch")
	}
	if weight < 1 || price < 1 {
		return 0, 0, errors.New("price or weight can't be negetive or zero")
	}
	return price + 1, weight, nil
}
