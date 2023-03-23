package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Handler interface {
	Save(context.Context, StorableObject) error
	First(context.Context, StorableObject) (StorableObject, error)
	Delete(context.Context, StorableObject) error
}

type StorableObject interface {
	ID() string
	SetID(string)
}

type InMemoryDB struct {
	collection map[string]StorableObject
}

func NewInMemoryDB() Handler {
	return &InMemoryDB{collection: make(map[string]StorableObject)}
}

func (sql *InMemoryDB) Save(ctx context.Context, val StorableObject) error {
	val.SetID(uuid.New().String())
	sql.collection[val.ID()] = val
	return nil
}

func (sql *InMemoryDB) First(ctx context.Context, val StorableObject) (StorableObject, error) {
	res, ok := sql.collection[val.ID()]
	if !ok {
		return nil, fmt.Errorf("record [%s] not found", val.ID())
	}

	return res, nil
}

func (sql *InMemoryDB) Delete(ctx context.Context, val StorableObject) error {
	delete(sql.collection, val.ID())
	return nil
}
