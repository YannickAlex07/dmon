package keiho

import "context"

type Monitor interface {
	Start(ctx context.Context, checkers []Checker, handlers []Handler, storage Storage) error
}
