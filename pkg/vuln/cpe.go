/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package vuln

import (
	"fmt"

	"github.com/facebookincubator/nvdtools/wfn"
)

// ValidateCPE returns an error if the CPE is invalid.
func ValidateCPE(cpe string) error {
	if cpe == "" {
		return fmt.Errorf("CPE must not be empty")
	}

	_, err := wfn.Parse(cpe)
	if err != nil {
		return fmt.Errorf("invalid CPE %q: %w", cpe, err)
	}

	return nil
}
