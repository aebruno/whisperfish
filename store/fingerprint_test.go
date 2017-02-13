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
	"testing"
)

// Test generating numeric fingerprints. This is a port of and all example data
// taken from the NumericFingerprintGeneratorTest from libsignal-protocol-java.
func TestNumericFingerprint(t *testing.T) {

	aliceKey := []byte{0x05, 0x06, 0x86, 0x3b, 0xc6, 0x6d, 0x02, 0xb4, 0x0d,
		0x27, 0xb8, 0xd4, 0x9c, 0xa7, 0xc0, 0x9e, 0x92, 0x39,
		0x23, 0x6f, 0x9d, 0x7d, 0x25, 0xd6, 0xfc, 0xca, 0x5c,
		0xe1, 0x3c, 0x70, 0x64, 0xd8, 0x68}
	bobKey := []byte{0x05, 0xf7, 0x81, 0xb6, 0xfb, 0x32, 0xfe, 0xd9, 0xba, 0x1c,
		0xf2, 0xde, 0x97, 0x8d, 0x4d, 0x5d, 0xa2, 0x8d, 0xc3, 0x40,
		0x46, 0xae, 0x81, 0x44, 0x02, 0xb5, 0xc0, 0xdb, 0xd9, 0x6f,
		0xda, 0x90, 0x7b}

	displayableFP := "300354477692869396892869876765458257569162576843440918079131"

	aliceFP := NumericFingerprint("+14152222222", aliceKey, "+14153333333", bobKey)
	bobFP := NumericFingerprint("+14153333333", bobKey, "+14152222222", aliceKey)

	if aliceFP != displayableFP {
		t.Errorf("Invalid numeric fingeprint for alice: got %s should be %s", aliceFP, displayableFP)
	}
	if bobFP != displayableFP {
		t.Errorf("Invalid numeric fingeprint for bob: got %s should be %s", bobFP, displayableFP)
	}
}
