package main

import "tolling/types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	return nil
}
