package repository

import (
	"context"

	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type distributorQueryRepository struct {
	db *gorm.DB
}

func NewDistributorQueryRepository(db *gorm.DB) repositories.DistributorQueryRepositoryInterface {
	return &distributorQueryRepository{
		db: db,
	}
}

func (r *distributorQueryRepository) FindByID(ctx context.Context, id uint) (*entities.Distributor, error) {
	var distributor entities.Distributor
	err := r.db.WithContext(ctx).First(&distributor, id).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

func (r *distributorQueryRepository) FindByEmail(ctx context.Context, email string) (*entities.Distributor, error) {
	var distributor entities.Distributor
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&distributor).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

func (r *distributorQueryRepository) FindAll(ctx context.Context) ([]*entities.Distributor, error) {
	var distributors []*entities.Distributor
	err := r.db.WithContext(ctx).Find(&distributors).Error
	return distributors, err
}