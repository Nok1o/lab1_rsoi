package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Nok1o/lab1_rsoi/backend/internal/domain"
)

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person domain.Person) (int64, error) {
	const query = `
		INSERT INTO persons (name, age, address, work)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, person.Name, person.Age, person.Address, person.Work).Scan(&id)
	return id, err
}

func (r *PersonRepository) GetByID(ctx context.Context, id int64) (domain.Person, error) {
	const query = `SELECT id, name, age, address, work FROM persons WHERE id = $1`

	var person domain.Person
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Age,
		&person.Address,
		&person.Work,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Person{}, domain.ErrPersonNotFound
	}
	return person, err
}

func (r *PersonRepository) List(ctx context.Context) ([]domain.Person, error) {
	const query = `SELECT id, name, age, address, work FROM persons ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := make([]domain.Person, 0)
	for rows.Next() {
		var person domain.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *PersonRepository) Update(ctx context.Context, person domain.Person) error {
	const query = `
		UPDATE persons
		SET name = $2, age = $3, address = $4, work = $5
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, person.ID, person.Name, person.Age, person.Address, person.Work)
	if err != nil {
		return err
	}
	return notFoundIfNoRows(result)
}

func (r *PersonRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM persons WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return notFoundIfNoRows(result)
}

func notFoundIfNoRows(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrPersonNotFound
	}
	return nil
}
