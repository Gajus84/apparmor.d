// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package cfg

import "github.com/arduino/go-paths-helper"

var (
	// Root is the root directory for the build
	Root *paths.Path = paths.New(".build")

	// RootApparmord is the final built apparmor.d directory
	RootApparmord *paths.Path = Root.Join("apparmor.d")

	// DistDir is the directory where the distribution specific files are stored
	DistDir *paths.Path = paths.New("dists")

	// FlagDir is the directory where the flags are stored
	FlagDir *paths.Path = DistDir.Join("flags")

	// SystemdDir is the directory where the systemd drop-in files are stored
	SystemdDir *paths.Path = paths.New("systemd")

	// Either or not overwrite some upstreamed profile
	Overwrite bool = false
)