package simpleretr

import (
	b64 "encoding/base64"
)

type MockDataStore interface {
	HasCid(cid string) (bool, error)
	GetBytes(cid string, offset int, count int) ([]byte, error)
}

type mockDataStore struct {
	mockCidBytes []byte
}

var _ MockDataStore = (*mockDataStore)(nil)

func NewMockDataStore() MockDataStore {
	dec, _ := b64.StdEncoding.DecodeString(CidBytesBase64)
	return &mockDataStore{
		mockCidBytes: dec,
	}
}

func (ds *mockDataStore) HasCid(cid string) (bool, error) {
	return (cid == "bafykbzacebcklmjetdwu2gg5svpqllfs37p3nbcjzj2ciswpszajbnw2ddxzo"), nil
}

func (ds *mockDataStore) GetBytes(cid string, offset int, count int) ([]byte, error) {
	bytes := []byte("hello")
	return bytes, nil
}
