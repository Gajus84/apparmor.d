// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Dbus directive
//
// Example of supported directive:
// #aa:dbus own bus=session name=org.freedesktop.FileManager1
// #aa:dbus talk name=org.freedesktop.login1 label=systemd-logind
//
// See https://apparmor.pujol.io/development/dbus/#dbus-directive
//

package directive

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
)

var defaultInterfaces = []string{
	"org.freedesktop.DBus.Properties",
	"org.freedesktop.DBus.ObjectManager",
}

type Dbus struct {
	DirectiveBase
}

func init() {
	Directives["dbus"] = &Dbus{
		DirectiveBase: DirectiveBase{
			message: "Dbus directive applied",
			usage: `#aa:dbus own bus=(system | session) name=<interface>
    #aa:dbus talk bus=(system | session) name=<interface> label=<profile_name>`,
		},
	}
}

func setInterfaces(rules map[string]string) []string {
	interfaces := []string{rules["name"]}
	if _, present := rules["interface"]; present {
		interfaces = append(interfaces, rules["interface"])
	}
	interfaces = append(interfaces, defaultInterfaces...)
	return interfaces
}

func (d Dbus) Apply(opt *Option, profile string) string {
	var p *aa.AppArmorProfile

	action := d.sanityCheck(opt)
	switch action {
	case "own":
		p = d.own(opt.ArgMap)
	case "talk":
		p = d.talk(opt.ArgMap)
	}

	generatedDbus := p.String()
	lenDbus := len(generatedDbus)
	generatedDbus = generatedDbus[:lenDbus-1]
	profile = strings.Replace(profile, opt.Raw, generatedDbus, -1)
	return profile
}

func (d Dbus) sanityCheck(opt *Option) string {
	if len(opt.ArgList) < 1 {
		panic(fmt.Sprintf("Unknown dbus action: %s in %s", opt.Name, opt.File))
	}
	action := opt.ArgList[0]
	if action != "own" && action != "talk" {
		panic(fmt.Sprintf("Unknown dbus action: %s in %s", opt.Name, opt.File))
	}

	if _, present := opt.ArgMap["name"]; !present {
		panic(fmt.Sprintf("Missing name for 'dbus: %s' in %s", action, opt.File))
	}
	if _, present := opt.ArgMap["bus"]; !present {
		panic(fmt.Sprintf("Missing bus for '%s' in %s", opt.ArgMap["name"], opt.File))
	}
	if _, present := opt.ArgMap["label"]; !present && action == "talk" {
		panic(fmt.Sprintf("Missing label for '%s' in %s", opt.ArgMap["name"], opt.File))
	}

	// Set default values
	if _, present := opt.ArgMap["path"]; !present {
		opt.ArgMap["path"] = "/" + strings.Replace(opt.ArgMap["name"], ".", "/", -1) + "{,/**}"
	}
	opt.ArgMap["name"] += "{,.*}"
	return action
}

func (d Dbus) own(rules map[string]string) *aa.AppArmorProfile {
	interfaces := setInterfaces(rules)
	p := &aa.AppArmorProfile{}
	p.Rules = append(p.Rules, &aa.Dbus{
		Access: "bind", Bus: rules["bus"], Name: rules["name"],
	})
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "receive",
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			Name:      `":1.@{int}"`,
		})
	}
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "send",
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			Name:      `"{:1.@{int},org.freedesktop.DBus}"`,
		})
	}
	p.Rules = append(p.Rules, &aa.Dbus{
		Access:    "receive",
		Bus:       rules["bus"],
		Path:      rules["path"],
		Interface: "org.freedesktop.DBus.Introspectable",
		Member:    "Introspect",
		Name:      `":1.@{int}"`,
	})
	return p
}

func (d Dbus) talk(rules map[string]string) *aa.AppArmorProfile {
	interfaces := setInterfaces(rules)
	p := &aa.AppArmorProfile{}
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "send",
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			Name:      `"{:1.@{int},` + rules["name"] + `}"`,
			Label:     rules["label"],
		})
	}
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "receive",
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			Name:      `"{:1.@{int},` + rules["name"] + `}"`,
			Label:     rules["label"],
		})
	}
	return p
}