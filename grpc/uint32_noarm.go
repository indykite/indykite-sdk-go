// Copyright (c) 2024 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !arm && !arm64
// +build !arm,!arm64

package grpc

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// Uint32 converts an integer of any type to uint32 safely.
// It returns an error if the value is out of the uint32 range.
func Uint32[T constraints.Integer](i T) (uint32, error) {
	// Check if i is negative (only relevant for signed integers).
	var zero T
	if i < zero {
		return 0, fmt.Errorf("uint32 out of range %d", i)
	}

	// For both signed and unsigned integers, ensure the value does not exceed MaxUint32.
	if uint64(i) > math.MaxUint32 {
		return 0, fmt.Errorf("uint32 out of range %d", i)
	}

	// Safe to convert to uint32 now.
	return uint32(i), nil
}
