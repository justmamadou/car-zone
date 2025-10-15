package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	Year      string    `json:"year"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	FuelType  string    `json:"fuel_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CarRequest struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	Year      string    `json:"year"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	FuelType  string    `json:"fuel_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func validateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("year cannot be empty")
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year must be a valid number")
	}
	currentYear := time.Now().Year()
	if yearInt < 1886 || yearInt > currentYear {
		return errors.New("year must be between 1886 and the current year")
	}
	return nil

}

func validateBrand(brand string) error {
	if brand == "" {
		return errors.New("brand cannot be empty")
	}
	return nil
}

func validateFuelType(fuelType string) error {
	validateFuelTypes := []string{"Petrol", "Diesel", "Electric", "Hybrid"}
	for _, v := range validateFuelTypes {
		if fuelType == v {
			return nil
		}
	}
	return errors.New("invalid fuel type, must be one of Petrol, Diesel, Electric, Hybrid")
}

func validateEngine(engine Engine) error {
	if engine.EngineID == uuid.Nil {
		return errors.New("engine_id cannot be empty")
	}
	if engine.Displacement <= 0 {
		return errors.New("displacement must be a positive number")
	}
	if engine.NoOfCylinders <= 0 {
		return errors.New("no_of_cylinders must be a positive number")
	}
	if engine.CarRange < 0 {
		return errors.New("car_range cannot be negative")
	}
	return nil
}

func validatePrice(price float64) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

func validateCarRequest(car CarRequest) error {
	if err := validateName(car.Name); err != nil {
		return err
	}
	if err := validateYear(car.Year); err != nil {
		return err
	}
	if err := validateBrand(car.Brand); err != nil {
		return err
	}
	if err := validateFuelType(car.FuelType); err != nil {
		return err
	}
	if err := validateEngine(car.Engine); err != nil {
		return err
	}
	if err := validatePrice(car.Price); err != nil {
		return err
	}
	return nil
}
