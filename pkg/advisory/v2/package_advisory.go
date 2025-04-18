/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

// PackageAdvisory is an Advisory that includes the package name for that
// advisory. (The Advisory type does not include the package's name.)
type PackageAdvisory struct {
	PackageName string `yaml:"packageName" json:"packageName"`
	Advisory
}

// IsZero returns true if the package advisory has no data.
func (pa PackageAdvisory) IsZero() bool {
	return pa.PackageName == "" && pa.Advisory.IsZero()
}
