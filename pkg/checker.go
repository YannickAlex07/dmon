package gmon

import "context"

type Checker interface {
	Check(ctx context.Context, handlers []Handler) error
}
