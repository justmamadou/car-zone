package car

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/justmamadou/car-zone/models"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) GetCarByID(ctx context.Context, carID string) (models.Car, error) {
	var car models.Car
	query := `SELECT c.id, c.name, c.brand, c.year, c.engine_id, c.price, c.fuel_type, c.created_at, c.updated_at, e.displacement, e.no_of_cylinders, e.car_range
	 FROM car c LEFT JOIN engine e ON c.engine_id = e.id
	 WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, carID)
	err := row.Scan(&car.ID, &car.Name, &car.Brand, &car.Year, &car.Engine.EngineID, &car.Price, &car.FuelType, &car.CreatedAt, &car.UpdatedAt, &car.Engine.Displacement, &car.Engine.NoOfCylinders, &car.Engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return car, nil
		}
		return car, err
	}
	return car, nil

}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {

	var car []models.Car
	var query string

	if isEngine {
		query = `SELECT c.id, c.name, c.brand, c.year, c.engine_id, c.price, c.fuel_type, c.created_at, c.updated_at, e.displacement, e.no_of_cylinders, e.car_range
		FROM car c 
		LEFT JOIN engine e ON c.engine_id = e.id
		WHERE c.brand = $1`
	} else {
		query = `SELECT id, name, brand, year, engine_id, price, fuel_type, created_at, updated_at
		FROM car 
		WHERE brand = $1`
	}

	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Car
		if isEngine {
			var engine models.Engine
			err = rows.Scan(&c.ID, &c.Name, &c.Brand, &c.Year, &c.Engine.EngineID, &c.Price, &c.FuelType, &c.CreatedAt, &c.UpdatedAt, &c.Engine.Displacement, &c.Engine.NoOfCylinders, &c.Engine.CarRange)
			if err != nil {
				return nil, err
			}
			c.Engine = engine

		} else {
			err = rows.Scan(&c.ID, &c.Name, &c.Brand, &c.Year, &c.Engine.EngineID, &c.Price, &c.FuelType, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				return nil, err

			}
		}
		car = append(car, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return car, nil

}

func (s Store) CreateCar(ctx context.Context, car *models.CarRequest) (models.Car, error) {

	var createdCar models.Car
	var engineID uuid.UUID

	err := s.db.QueryRowContext(ctx, "SELECT id FROM engine WHERE id=$1", car.Engine.EngineID).Scan(&engineID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return createdCar, errors.New("engine not found")
		}
		return createdCar, err
	}

	carID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	newCar := models.Car{
		ID:        carID,
		Name:      car.Name,
		Brand:     car.Brand,
		Year:      car.Year,
		Engine:    car.Engine,
		Price:     car.Price,
		FuelType:  car.FuelType,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	//Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return createdCar, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	insertQuery := `INSERT INTO car (id, name, brand, year, engine_id, price, fuel_type, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, name, brand, year, engine_id, price, fuel_type, created_at, updated_at`

	err = tx.QueryRowContext(ctx, insertQuery, newCar.ID, newCar.Name, newCar.Brand, newCar.Year, newCar.Engine.EngineID,
		newCar.Price, newCar.FuelType, newCar.CreatedAt, newCar.UpdatedAt).Scan(&createdCar.ID, &createdCar.Name, &createdCar.Brand,
		&createdCar.Year, &createdCar.Engine.EngineID, &createdCar.Price, &createdCar.FuelType, &createdCar.CreatedAt, &createdCar.UpdatedAt)
	if err != nil {
		return createdCar, err
	}
	return createdCar, nil
}

func (s Store) UpdateCar(ctx context.Context, id string, car *models.CarRequest) (models.Car, error) {

	carID, err := uuid.Parse(id)
	if err != nil {
		return models.Car{}, errors.New("invalid Car ID format")
	}

	var updatedCar models.Car
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedCar, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	updateQuery := `UPDATE car SET name=$2, brand=$3, year=$4, engine_id=$5, price=$6, fuel_type=$7, updated_at=$8
	WHERE id=$1
	RETURNING id, name, brand, year, engine_id, price, fuel_type, created_at, updated_at`

	err = tx.QueryRowContext(ctx, updateQuery, carID, car.Name, car.Brand, car.Year, car.Engine.EngineID,
		car.Price, car.FuelType, time.Now()).Scan(&updatedCar.ID, &updatedCar.Name, &updatedCar.Brand,
		&updatedCar.Year, &updatedCar.Engine.EngineID, &updatedCar.Price, &updatedCar.FuelType, &updatedCar.CreatedAt, &updatedCar.UpdatedAt)
	if err != nil {
		return updatedCar, err
	}
	return updatedCar, nil

}

func (s Store) DeleteCar(ctx context.Context, carID string) (models.Car, error) {
	var deletedCar models.Car
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return deletedCar, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	checkQuery := `SELECT id, name, brand, year, engine_id, price, fuel_type, created_at, updated_at
	FROM car
	WHERE id=$1`

	err = tx.QueryRowContext(ctx, checkQuery, carID).Scan(&deletedCar.ID, &deletedCar.Name, &deletedCar.Brand,
		&deletedCar.Year, &deletedCar.Engine.EngineID, &deletedCar.Price, &deletedCar.FuelType, &deletedCar.CreatedAt, &deletedCar.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deletedCar, errors.New("car not found")
		}
		return deletedCar, err
	}

	deleteQuery := `DELETE FROM car WHERE id=$1`

	result, err := tx.ExecContext(ctx, deleteQuery, carID)
	if err != nil {
		return deletedCar, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return deletedCar, err
	}

	if rowsAffected == 0 {
		return deletedCar, errors.New("car not found")
	}

	return deletedCar, nil
}
