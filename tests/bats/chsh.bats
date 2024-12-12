#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "chsh: [l]ist available shells" {
    chsh --list-shells || true
}

@test "chsh: Set a specific login [s]hell for the current user" {
    echo "$PASSWORD" | chsh --shell /usr/bin/bash
}

# bats test_tags=chsh
@test "chsh: Set a login [s]hell for a specific user" {
    sudo chsh --shell /usr/bin/sh root
}