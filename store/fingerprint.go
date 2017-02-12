// Copyright 2016 Andrew E. Bruno
//
// This file is part of Whisperfish.
//
// Whisperfish is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Whisperfish is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Whisperfish.  If not, see <http://www.gnu.org/licenses/>.

package store

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	// Signal fingerprint version
	FingerprintVersion = 0

	// The number of internal iterations to perform in the process of
	// generating a fingerprint
	Iterations = 5200
)

// Generate a displayable numeric fingerprint from the client's "stable"
// identifier localStableIdentifier, the client's identity key
// localIdentityKey, the remort party's stable identifier
// remoteStableIdentifier, and the remote party's identity key. Returns a
// unique fingerprint for the conversation. This is a port of
// libsignal-protocol-java NumericFingerprintGenerator.
func NumericFingerprint(localStableIdentifier string, localIdentityKey []byte, remoteStableIdentifier string, remoteIdentityKey []byte) string {
	localFP := createFingerprint([]byte(localStableIdentifier), localIdentityKey)
	remoteFP := createFingerprint([]byte(remoteStableIdentifier), remoteIdentityKey)

	if strings.Compare(localFP, remoteFP) <= 0 {
		return localFP + remoteFP
	} else {
		return remoteFP + localFP
	}
}

func createFingerprint(stableIdentifier []byte, publicKey []byte) string {
	hasher := sha512.New()

	hash := make([]byte, 2)
	binary.LittleEndian.PutUint16(hash, FingerprintVersion)
	hash = append(hash, publicKey[:]...)
	hash = append(hash, stableIdentifier[:]...)

	for i := 0; i < Iterations; i++ {
		hasher.Write(hash)
		hasher.Write(publicKey)
		hash = hasher.Sum(nil)
		hasher.Reset()
	}

	var buffer bytes.Buffer

	buffer.WriteString(encodeChunk(hash, 0))
	buffer.WriteString(encodeChunk(hash, 5))
	buffer.WriteString(encodeChunk(hash, 10))
	buffer.WriteString(encodeChunk(hash, 15))
	buffer.WriteString(encodeChunk(hash, 20))
	buffer.WriteString(encodeChunk(hash, 25))

	return buffer.String()
}

func encodeChunk(fp []byte, offset int) string {
	n := ((int64(fp[offset]) & 0xff) << 32) |
		((int64(fp[offset+1]) & 0xff) << 24) |
		((int64(fp[offset+2]) & 0xff) << 16) |
		((int64(fp[offset+3]) & 0xff) << 8) |
		(int64(fp[offset+4]) & 0xff)

	return fmt.Sprintf("%05d", n%100000)
}
