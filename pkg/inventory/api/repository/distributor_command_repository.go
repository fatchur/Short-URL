package repository

import (
	"context"

	"short-url/domains/entities"
	"short-url/domains/repositories"

	"gorm.io/gorm"
)

type distributorCommandRepository struct {
	db *gorm.DB
}

func NewDistributorCommandRepository(db *gorm.DB) repositories.DistributorCommandRepositoryInterface {
	return &distributorCommandRepository{
		db: db,
	}
}

func (r *distributorCommandRepository) Save(ctx context.Context, distributor *entities.Distributor) error {
	return r.db.WithContext(ctx).Create(distributor).Error
}

func (r *distributorCommandRepository) Update(ctx context.Context, distributor *entities.Distributor) error {
	return r.db.WithContext(ctx).Save(distributor).Error
}

func (r *distributorCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Distributor{}, id).Error
}