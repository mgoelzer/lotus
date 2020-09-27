package simpleretr

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
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
	fmt.Printf("[sretrieve] (GetBytes) Entering function with offset=%v, count=%v, cid='%v'\n", offset, count, cid)
	hasCid, _ := ds.HasCid(cid)
	fmt.Printf("[sretrieve] (GetBytes) hasCid='%v'\n", hasCid)
	if !hasCid {
		return []byte(nil), errors.New("[sretrieve] (GetBytes) Unrecognized cid")
	}

	fmt.Printf("[sretrieve] (GetBytes) size of ds.mockCidBytes ='%v'\n", len(ds.mockCidBytes))

	log.Infof("[sretrieve] (GetBytes) Starting to load a []byte with mock data\n")
	bytes := []byte(ds.mockCidBytes)
	log.Infof("[sretrieve] (GetBytes)  Done loading\n")

	fmt.Printf("[sretrieve] (GetBytes) First 8 bytes of Cid: %x %x %x %x %x %x %x %x\n", bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7])

	retBytes := bytes[offset : offset+count]
	//log.Infof
	fmt.Printf("[sretrieve] (GetBytes)  Returning bytes[%v:%v] with len %v bytes\n", offset, offset+count, len(retBytes))
	return retBytes, nil
}
