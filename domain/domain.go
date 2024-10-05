package domain

// LoadDomain loads the various structs to their global variables.
// Since some of them reference images from assets, it should be called after assets.MustLoadAssets.
func LoadDomain() {
	loadTerrains()
	loadDeposits()
	loadPlayerCharacters()
	loadTechnologies()
	loadResources()
}
