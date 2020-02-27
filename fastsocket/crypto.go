// Copyright 2019 Vinzenz Weist. All rights reserved.
// Use of this source code is risked by yourself.
// license that can be found in the LICENSE file.
package fastsocket

import (
	"crypto/sha256"
	"github.com/google/uuid"

)
// isUUID checks if it's a valid uuid string
func isUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
// generates a sha256 byte array from a string
func generateSHA256(str string) [32]byte {
	sha256 := sha256.Sum256([]byte(str))
	return sha256
}