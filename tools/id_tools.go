package tools

import (
	"math/rand"

	"github.com/grayscalecloud/kitexcommon/utils"
)

func GenerateShareNo() string {
	nodeID := rand.Intn(31) + 1
	sf := utils.NewSnowflake(int64(nodeID))
	id := sf.EncodeToShortID()
	return id
}
