package service

import "context"

type Service interface {
	GetLastAction(ctx context.Context) (string, error)
}
