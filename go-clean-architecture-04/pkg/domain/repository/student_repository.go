package repository

import (
	"context"

	"go-clean_architecture-api/pkg/domain/model"
)

// IHogeHoge represents interface of HogeHoge
type IStudentRepository interface {
	SelectAllStudents(ctx context.Context) (model.StudentSlice, error)
	SelectStudentByID(ctx context.Context, id int) (*model.Student, error)
}
