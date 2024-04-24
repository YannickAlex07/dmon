package inframon_test

import (
	"context"
	"errors"
	"time"

	inframon "github.com/yannickalex07/inframon/pkg"
)

// Source

type FakeSource struct {
	Notifications []inframon.Notification
}

func (f *FakeSource) Check(ctx context.Context, since time.Time) ([]inframon.Notification, error) {
	return f.Notifications, nil
}

// Destination

type FakeDestination struct {
	Notifications []inframon.Notification
}

func (f *FakeDestination) Handle(ctx context.Context, notification inframon.Notification) error {
	f.Notifications = append(f.Notifications, notification)
	return nil
}

// State

type FakeState struct {
	Storage map[string]interface{}
}

func (f *FakeState) Store(ctx context.Context, key string, value interface{}, shouldExpire bool) error {
	f.Storage[key] = value
	return nil
}

func (f *FakeState) Get(ctx context.Context, key string) (interface{}, error) {
	if val, ok := f.Storage[key]; ok {
		return val, nil
	}

	return nil, errors.New("key not found")
}

func (f *FakeState) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := f.Storage[key]
	return ok, nil
}
