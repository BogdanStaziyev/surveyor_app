package database

import (
	"fmt"
	"github.com/upper/db/v4"
	"log"
	"startUp/internal/domain"
	"time"
)

const CoordinateTable = "coordinate"

type coordinate struct {
	Id          int64      `db:"id,omitempty"`
	MT          int64      `db:"mt"`
	Axis        string     `db:"axis"`
	Horizon     string     `db:"horizon"`
	X           float64    `db:"x"`
	Y           float64    `db:"y"`
	UserId      int64      `db:"user_id"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date, omitempty"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type CoordinateRepository interface {
	AddCoordinate(coordinate *domain.Coordinate) (*domain.Coordinate, error)
	UpdateCoordinate(coordinate *domain.Coordinate) error
	DeleteCoordinate(id int64) error
	FindAll() ([]domain.Coordinate, error)
	FindOne(id int64) (*domain.Coordinate, error)
	InverseTask(firstId, secondId int64) (string, error, *domain.Coordinate, *domain.Coordinate)
}

type repository struct {
	coll db.Collection
}

func NewRepository(dbSession *db.Session) CoordinateRepository {
	return &repository{
		coll: (*dbSession).Collection("coordinate"),
	}
}

func (r *repository) AddCoordinate(coordinate *domain.Coordinate) (*domain.Coordinate, error) {
	coordinates := mapCoordinateDbModel(coordinate)
	coordinates.CreatedDate = time.Now()
	err := r.coll.InsertReturning(coordinates)
	if err != nil {
		return nil, fmt.Errorf("CoordinaterepositoryCreate: %w", err)
	}

	return mapCoordinateDbModelToDomain(coordinates), nil
}

func (r *repository) UpdateCoordinate(coordinate *domain.Coordinate) error {
	coordinates := mapCoordinateDbModel(coordinate)

	coordinates.UpdatedDate = time.Now()
	err := r.coll.Find(coordinates.Id).Update(coordinates)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (r *repository) DeleteCoordinate(id int64) error {

	err := r.coll.Find("id", id).Delete()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (r *repository) FindAll() ([]domain.Coordinate, error) {
	var coordinates []coordinate

	err := r.coll.Find().All(&coordinates)
	if err != nil {
		log.Fatal("coordinate.Find: ", err)
	}
	return mapCoordinateDbModelToDomainCollection(coordinates), nil
}

func (r *repository) FindOne(id int64) (*domain.Coordinate, error) {
	var coordinates coordinate

	err := r.coll.Find("id", id).One(&coordinates)
	if err != nil {
		return nil, fmt.Errorf("coordinateRepository FindOne: %w", err)
	}
	return mapCoordinateDbModelToDomain(&coordinates), nil
}

func (r *repository) InverseTask(firstId, secondId int64) (string, error, *domain.Coordinate, *domain.Coordinate) {
	var coordinateOne, coordinateTwo coordinate
	firstErr := r.coll.Find("id", firstId).One(&coordinateOne)
	if firstErr != nil {
		log.Fatal("repository Invert first: ", firstErr)
	}
	secondErr := r.coll.Find("id", secondId).One(&coordinateTwo)
	if secondErr != nil {
		log.Fatal("repository Invert second: ", secondErr)
	}
	return "Результат обчислення зворотньої геодезичної задачі: ", nil, mapCoordinateDbModelToDomain(&coordinateOne), mapCoordinateDbModelToDomain(&coordinateTwo)
}

func mapCoordinateDbModelToDomain(coordinate *coordinate) *domain.Coordinate {
	return &domain.Coordinate{
		Id:          coordinate.Id,
		MT:          coordinate.MT,
		Axis:        coordinate.Axis,
		Horizon:     coordinate.Horizon,
		X:           coordinate.X,
		Y:           coordinate.Y,
		UserID:      coordinate.UserId,
		CreatedDate: coordinate.CreatedDate,
		UpdatedDate: coordinate.UpdatedDate,
		DeletedDate: getTimeFromTimePtr(coordinate.DeletedDate),
	}
}

func mapCoordinateDbModelToDomainCollection(coordinate []coordinate) []domain.Coordinate {
	var result []domain.Coordinate
	for _, c := range coordinate {
		newCoordinate := mapCoordinateDbModelToDomain(&c)
		result = append(result, *newCoordinate)
	}
	return result
}

func mapCoordinateDbModel(coord *domain.Coordinate) *coordinate {
	return &coordinate{
		Id:          coord.Id,
		MT:          coord.MT,
		Axis:        coord.Axis,
		Horizon:     coord.Horizon,
		X:           coord.X,
		Y:           coord.Y,
		UserId:      coord.UserID,
		CreatedDate: coord.CreatedDate,
		UpdatedDate: coord.UpdatedDate,
		DeletedDate: getTimePtrFromTime(coord.DeletedDate),
	}
}
