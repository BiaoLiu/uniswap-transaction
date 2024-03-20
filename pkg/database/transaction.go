package database

import "context"

type Transaction interface {
	InTx(context.Context, string, func(ctx context.Context) error) error
}
