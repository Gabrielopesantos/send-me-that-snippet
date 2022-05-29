package paste

import (
	"context"
	"github.com/gabrielopesantos/smts/internal/model"
)

type Repository interface {
	Insert(context.Context, *model.Paste) error
	Delete(context.Context, string) (*model.Paste, error)
	Get(context.Context, string) (*model.Paste, error)
	Update(context.Context, string, *model.Paste) error
	Filter(context.Context, *model.Paste) ([]model.Paste, error)
}
