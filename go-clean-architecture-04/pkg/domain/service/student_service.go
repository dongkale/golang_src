package service

import (
	"context"

	"go-clean_architecture-api/pkg/domain/model"
	"go-clean_architecture-api/pkg/domain/repository"
)

// interface
type IStudentService interface {
	FindAllStudents(ctx context.Context) (model.StudentSlice, error)
	FindStudentByID(ctx context.Context, id int) (*model.Student, error)
}

// struct that meets interface
type studentService struct {
	repo repository.IStudentRepository
}

// constructor
func NewStudentService(sr repository.IStudentRepository) IStudentService {
	return &studentService{
		repo: sr,
	}
}

// implement methods according to interface
func (ss *studentService) FindAllStudents(ctx context.Context) (model.StudentSlice, error) {
	return ss.repo.SelectAllStudents(ctx)
}

func (ss *studentService) FindStudentByID(ctx context.Context, id int) (*model.Student, error) {
	return ss.repo.SelectStudentByID(ctx, id)
}
