//go:build !ios && !android

package storage

import "github.com/kirsle/configdir"

func storageDirectory() (string, error) {
	return configdir.LocalConfig("kingdoms"), nil
}
