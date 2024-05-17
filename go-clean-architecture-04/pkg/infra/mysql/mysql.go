package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"go-clean_architecture-api/pkg/domain/model"
	"go-clean_architecture-api/pkg/domain/repository"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type studentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) repository.IStudentRepository {
	return &studentRepository{
		DB: db,
	}
}

func (sr *studentRepository) SelectAllStudents(ctx context.Context) (model.StudentSlice, error) {
	// concrete DB operation
	return model.Students().All(ctx, sr.DB)
}

func (sr *studentRepository) SelectStudentByID(ctx context.Context, studentID int) (*model.Student, error) {
	// concrete DB operation
	whereID := fmt.Sprintf("%s = ?", model.StudentColumns.ID)
	return model.Students(
		qm.Where(whereID, studentID),
	).One(ctx, sr.DB)
}
