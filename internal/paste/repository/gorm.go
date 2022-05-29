package paste

import (
	"context"
	"errors"
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/internal/paste"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
	//logger *logger.Logger
}

func NewGormRepository(db *gorm.DB) paste.Repository {
	return &GormRepository{
		db: db,
		//logger: logger,
	}
}

func (gr *GormRepository) Insert(ctx context.Context, paste *model.Paste) error {
	_, span := otel.Tracer("main-service").Start(ctx, "Insert.DB")
	defer span.End()

	return gr.db.Create(&paste).Error
}

// Not returning deleted paste
func (gr *GormRepository) Delete(ctx context.Context, pId string) (*model.Paste, error) {
	_, span := otel.Tracer("main-service").Start(ctx, "Delete.DB")
	defer span.End()

	p := model.Paste{}
	//result := gr.db.Where("id = ?", pId).Delete(&p)
	result := gr.db.Delete(&p, "id = ?", pId)
	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}

	return &p, result.Error
}

func (gr *GormRepository) Get(ctx context.Context, pId string) (*model.Paste, error) {
	_, span := otel.Tracer("main-service").Start(ctx, "Get.DB")
	defer span.End()

	p := model.Paste{}
	err := gr.db.First(&p, "id  = ?", pId).Error

	return &p, err
}

func (gr *GormRepository) Update(ctx context.Context, pasteId string, paste *model.Paste) error {
	_, span := otel.Tracer("main-service").Start(ctx, "Update.DB")
	defer span.End()

	return gr.db.Where("id = ?", pasteId).Updates(paste).Error
}

func (gr *GormRepository) Filter(ctx context.Context, filter *model.Paste) ([]model.Paste, error) {
	_, span := otel.Tracer("main-service").Start(ctx, "Filter.DB")
	defer span.End()

	var results []model.Paste
	// This filter is currently broken due to the "expired" args
	err := gr.db.Where(&filter, "expired").Find(&results).Error

	return results, err
}
