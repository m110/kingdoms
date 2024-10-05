//go:build android

package storage

func storageDirectory() (string, error) {
	// TODO Probably should not be hardcoded
	return "/data/data/m110.games.kingdoms", nil
}
