//go:build !nogames
// +build !nogames

package main

func init() {
	registeredFeatures["games"] = func() error {
		return nil
	}
}
