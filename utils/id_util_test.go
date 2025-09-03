package utils

import (
	"fmt"
	"testing"
)

func Test_encodeToShortID(t *testing.T) {
	sf := NewSnowflake(1)
	_ = sf.NextID()
	shortID := sf.EncodeToShortID()

	fmt.Println(shortID)
}
