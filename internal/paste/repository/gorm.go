package paste

import (
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/internal/paste"
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

func (gr *GormRepository) Insert(paste *model.Paste) error {
	return gr.db.Create(&paste).Error
}

func (gr *GormRepository) Delete(pId string) (*model.Paste, error) {
	p := model.Paste{}
	err := gr.db.Delete(&p, "id = ?", pId).Error

	return &p, err
}

func (gr *GormRepository) Get(pId string) (*model.Paste, error) {
	p := model.Paste{}
	err := gr.db.First(&p, "id  = ?", pId).Error

	return &p, err
}

func (gr *GormRepository) Update(pasteId string, paste *model.Paste) error {
	return gr.db.Where("id = ?", pasteId).Updates(paste).Error
}

func (gr *GormRepository) Filter(filter *model.Paste) ([]model.Paste, error) {
	var results []model.Paste
	// This filter is currently boken due to the "expired" args
	err := gr.db.Where(&filter, "expired").Find(&results).Error

	return results, err
}
