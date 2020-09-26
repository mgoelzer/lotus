package simpleretr

import (
	b64 "encoding/base64"
	"errors"
)

type MockDataStore interface {
	HasCid(cid string) (bool, error)
	GetBytes(cid string, offset int64, count int64) ([]byte, error)
}

type mockDataStore struct {
	mockCidBytes []byte
}

// static type check that mockDataStore satisfies interface MockDataStore
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

func (ds *mockDataStore) GetBytes(cid string, offset int64, count int64) ([]byte, error) {
	hasCid, _ := ds.HasCid(cid)
	if !hasCid {
		return []byte(nil), errors.New("[sretrieve] (GetBytes) Unrecognized cid")
	}

	log.Infof("[sretrieve] (GetBytes) Starting to load a []byte with mock data\n")
	bytes := []byte(ds.mockCidBytes)
	log.Infof("[sretrieve] (GetBytes)  Done loading\n")

	log.Infof("[sretrieve] (GetBytes)  Returning bytes[%v:%v]\n", offset, offset+count)
	return bytes[offset : offset+count], nil
}
