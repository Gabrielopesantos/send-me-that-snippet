package paste

import (
	"github.com/gabrielopesantos/smts/internal/model"
)

type Repository interface {
	Insert(*model.Paste) error
	Delete(string) (*model.Paste, error)
	Get(string) (*model.Paste, error)
	Update(string, *model.Paste) error
	Filter(*model.Paste) ([]model.Paste, error)
}
