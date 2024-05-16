package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joisandresky/go-echo-mysql-boilerplate/database"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/entity"
)

type HumanRepository interface {
	GetAll(ctx context.Context) ([]entity.Human, error)
	Show(ctx context.Context, id int) (entity.Human, error)
	Store(ctx context.Context, human entity.Human) error
	Update(ctx context.Context, id int, human entity.Human) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
}

type humanRepository struct {
	c *sqlx.DB
}

func NewHumanRepository(conn database.DatabaseProviderConnection) HumanRepository {
	return &humanRepository{conn.Db}
}

func (r *humanRepository) GetAll(ctx context.Context) ([]entity.Human, error) {
	var humans []entity.Human
	rows, err := r.c.QueryxContext(ctx, "SELECT * FROM humans")
	if err != nil {
		return nil, err
	}

	var human entity.Human
	for rows.Next() {
		err := rows.StructScan(&human)
		if err != nil {
			log.Println("FAILED_TO_DECODE_STRUCT_FROM_ROW", err)
		}

		humans = append(humans, human)
	}
	defer rows.Close()

	return humans, nil
}

func (r *humanRepository) Show(ctx context.Context, id int) (entity.Human, error) {
	var human entity.Human
	row := r.c.QueryRowxContext(ctx, "SELECT * FROM humans WHERE id = ?", id)
	if err := row.StructScan(&human); err != nil {
		return human, err
	}

	return human, nil
}

func (r *humanRepository) Store(ctx context.Context, human entity.Human) error {
	stmt, err := r.c.PrepareContext(ctx, "INSERT INTO humans (name, race) VALUES(?,?)")
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, human.Name, human.Race)
	if err != nil {
		return err
	}
	inserted, _ := result.LastInsertId()
	log.Println("INSERTED", inserted)
	defer stmt.Close()

	return nil
}

func (r *humanRepository) Update(ctx context.Context, id int, human entity.Human) (int64, error) {
	stmt, err := r.c.PrepareContext(ctx, "UPDATE humans SET name = ?, race = ?  WHERE id = ?")
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, human.Name, human.Race, id)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	return result.RowsAffected()
}

func (r *humanRepository) Delete(ctx context.Context, id int) (int64, error) {
	stmt, err := r.c.PrepareContext(ctx, "DELETE FROM humans WHERE id = ?")
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	return result.RowsAffected()
}
