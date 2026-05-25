package artifact

import (
	"crypto/sha256"
	"fmt"
)

func calculateChecksum(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
