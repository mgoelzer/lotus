package simpleretr

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirst16Bytes(t *testing.T) {
	cid := "bafykbzacebcklmjetdwu2gg5svpqllfs37p3nbcjzj2ciswpszajbnw2ddxzo"
	N := int64(16)
	Offset := int64(0)

	var bytes []byte
	fmt.Printf("[sretrieve] (TestGetFirst16Bytes) len(bytes)=%v before GetBytes\n", len(bytes))
	bytes, err := ds.GetBytes(cid, Offset, N)
	if err != nil {
		t.Errorf("[sretrieve] (TestGetFirst16Bytes) GetBytes failed with '%v'\n", err)
	}
	fmt.Printf("[sretrieve] (TestGetFirst16Bytes) len(bytes)=%v after GetBytes\n", len(bytes))

	bytesHexStr := hex.EncodeToString(bytes)
	fmt.Printf("[sretrieve] (TestGetFirst16Bytes) : '%v'\n", bytesHexStr)

	assert.Equal(t, "2fcc87b1eb83fddf6937386d7df2469b", bytesHexStr, "Bytes should be the same")
}

func TestGetSecondWord(t *testing.T) {
	cid := "bafykbzacebcklmjetdwu2gg5svpqllfs37p3nbcjzj2ciswpszajbnw2ddxzo"
	N := int64(8)
	Offset := int64(8)

	var bytes []byte
	bytes, err := ds.GetBytes(cid, Offset, N)
	if err != nil {
		t.Errorf("[sretrieve] (TestGetSecondWord) GetBytes failed with '%v'\n", err)
	}
	fmt.Printf("[sretrieve] (TestGetSecondWord) len(bytes)=%v\n", len(bytes))

	bytesHexStr := hex.EncodeToString(bytes)
	fmt.Printf("[sretrieve] (TestGetSecondWord) second word = '%v'\n", bytesHexStr)

	assert.Equal(t, "6937386d7df2469b", bytesHexStr, "Bytes should be the same")
}

func TestGetFirst8BytesAfter1MiBOffset(t *testing.T) {
	cid := "bafykbzacebcklmjetdwu2gg5svpqllfs37p3nbcjzj2ciswpszajbnw2ddxzo"
	N := int64(8)
	Offset := int64(1048576)
	fmt.Printf("[sretrieve] (TestGetFirst8BytesAfter1MiBOffset) offset=%v, N=%v, cid=%v\n", Offset, N, cid)

	var bytes []byte
	bytes, err := ds.GetBytes(cid, Offset, N)
	if err != nil {
		t.Errorf("[sretrieve] (TestGetFirst8BytesAfter1MiBOffset) GetBytes failed with '%v'\n", err)
	}
	fmt.Printf("[sretrieve] (TestGetFirst8BytesAfter1MiBOffset) len(bytes)=%v\n", len(bytes))
	assert.Equal(t, int64(N), int64(len(bytes)), "Request size and retrieved size must match (GetBytes error probably)")

	bytesHexStr := hex.EncodeToString(bytes)
	fmt.Printf("[sretrieve] (TestGetFirst8BytesAfter1MiBOffset) first word after 1MiB skipped = '%v'\n", bytesHexStr)

	assert.Equal(t, "e4e6bdad8c09e7fd", bytesHexStr, "Bytes should be the same")
}

func TestGetLastByte(t *testing.T) {
	cid := "bafykbzacebcklmjetdwu2gg5svpqllfs37p3nbcjzj2ciswpszajbnw2ddxzo"
	N := int64(1)
	Offset := int64(68157440 - 1)
	fmt.Printf("[sretrieve] (TestGetLastByte) offset=%v, N=%v, cid=%v\n", Offset, N, cid)

	var bytes []byte
	bytes, err := ds.GetBytes(cid, Offset, N)
	if err != nil {
		t.Errorf("[sretrieve] (TestGetLastByte) GetBytes failed with '%v'\n", err)
	}
	fmt.Printf("[sretrieve] (TestGetLastByte) Got back this number of bytes (expecting 1): %v\n", len(bytes))
	assert.Equal(t, int64(N), int64(len(bytes)), "Request size and retrieved size must match (GetBytes error probably)")

	assert.Equal(t, uint8(bytes[0]), uint8(0x25), "Last byte is 0x25 (=37 decimal)")
}
