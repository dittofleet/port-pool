// Package app is port-pool's description of itself for go-cli-kit, kept in
// one place so the config and state paths, the update cache and the release
// assets cannot drift apart.
package app

import clikit "github.com/dittofleet/go-cli-kit"

// Name is the install name, the repo name, and the directory port-pool
// keeps its files under.
const Name = "port-pool"

// New describes this build. The version is stamped into main at build
// time, so main passes it in.
func New(version string) clikit.App {
	return clikit.App{Name: Name, Version: version}
}
