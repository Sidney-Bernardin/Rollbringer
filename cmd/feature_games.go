//go:build !nopages
// +build !nopages

package main

func init() {
	registeredFeatures["games"] = func() error {
		return nil
	}
}
