package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineID      uuid.UUID `json:"engine_id"`
	Displacement  int64     `json:"displacement"`
	NoOfCylinders int64     `json:"no_of_cylinders"`
	CarRange      int64     `json:"car_range"`
}

type EngineRequest struct {
	Displacement  int64 `json:"displacement"`
	NoOfCylinders int64 `json:"no_of_cylinders"`
	CarRange      int64 `json:"car_range"`
}

func ValidateEngineRequest(engine EngineRequest) error {
	if err := validateDisplacement(engine.Displacement); err != nil {
		return err
	}
	if err := validateNoOfCylinders(engine.NoOfCylinders); err != nil {
		return err
	}
	if err := validateCarRange(engine.CarRange); err != nil {
		return err
	}
	return nil
}

func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be a positive number")
	}
	return nil
}

func validateNoOfCylinders(noOfCylinders int64) error {
	if noOfCylinders <= 0 {
		return errors.New("no_of_cylinders must be a positive number")
	}
	return nil
}

func validateCarRange(carRange int64) error {
	if carRange < 0 {
		return errors.New("car_range cannot be negative")
	}
	return nil
}
