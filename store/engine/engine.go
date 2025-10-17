package engine

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/justmamadou/car-zone/models"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) EngineStore {
	return EngineStore{db: db}
}

func (s EngineStore) GetEngineByID(ctx context.Context, engineID string) (models.Engine, error) {
	var engine models.Engine
	query := `SELECT id, displacement, no_of_cylinders, car_range
	 FROM engine
	 WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, engineID)
	err := row.Scan(&engine.EngineID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return engine, nil
		}
		return engine, err
	}
	return engine, nil
}

func (s EngineStore) CreateEngine(ctx context.Context, engineRequest *models.EngineRequest) (models.Engine, error) {
	var CreatedEngine models.Engine
	var engineID uuid.UUID

	engineID = uuid.New()

	newEngine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineRequest.Displacement,
		NoOfCylinders: engineRequest.NoOfCylinders,
		CarRange:      engineRequest.CarRange,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return CreatedEngine, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT INTO engine (id, displacement, no_of_cylinders, car_range)
	VALUES ($1, $2, $3, $4) RETURNING id, displacement, no_of_cylinders, car_range`
	
	err = tx.QueryRowContext(ctx, query, newEngine.EngineID, 
		newEngine.Displacement, newEngine.NoOfCylinders, newEngine.CarRange).Scan(&CreatedEngine.EngineID, 
			&CreatedEngine.Displacement, &CreatedEngine.NoOfCylinders, &CreatedEngine.CarRange)
	if err != nil {
		return CreatedEngine, err
	}
	return CreatedEngine, nil
}
func (s EngineStore) UpdateEngine(ctx context.Context, id string, engine *models.EngineRequest) (models.Engine, error) {
	engineID, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, errors.New("invalid engine ID format")
	}

	var updatedEngine models.Engine
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedEngine, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `UPDATE engine 
	SET displacement = $2, no_of_cylinders = $3, car_range = $4
	WHERE id = $1
	RETURNING id, displacement, no_of_cylinders, car_range`

	row := tx.QueryRowContext(ctx, query, engineID, engine.Displacement, engine.NoOfCylinders, engine.CarRange)
	err = row.Scan(&updatedEngine.EngineID, &updatedEngine.Displacement, &updatedEngine.NoOfCylinders, &updatedEngine.CarRange)
	if err != nil {
		return updatedEngine, err
	}
	return updatedEngine, nil
}

func (s EngineStore) DeleteEngine(ctx context.Context, engineID string) (models.Engine, error) {
	var deletedCar models.Engine
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

	checkQuery := `SELECT id, displacement, no_of_cylinders, car_range
	FROM engine
	WHERE id = $1`

	err = tx.QueryRowContext(ctx, checkQuery, engineID).Scan(&deletedCar.EngineID, &deletedCar.Displacement,
		&deletedCar.NoOfCylinders, &deletedCar.CarRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deletedCar, errors.New("engine not found")
		}
		return deletedCar, err
	}

	deleteQuery := `DELETE FROM engine WHERE id = $1`
	result, err := tx.ExecContext(ctx, deleteQuery, engineID)
	if err != nil {
		return deletedCar, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return deletedCar, err
	}
	if rowsAffected == 0 {
		return deletedCar, errors.New("no engine deleted")
	}
	return deletedCar, nil
}
