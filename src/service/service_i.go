package service

import "context"

type ServiceI interface {
	Start() error
	Shutdown(ctx context.Context) error
}
