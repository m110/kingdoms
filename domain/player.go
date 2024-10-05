package domain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/assets"
)

type PlayerCharacterType int

const (
	PlayerCharacterKing PlayerCharacterType = iota
	PlayerCharacterQueen
	PlayerCharacterPrincess
	PlayerCharacterMerchant
)

type PlayerCharacter struct {
	Name        string
	AddressedBy string
	Description string
	Image       *ebiten.Image
}

var PlayerCharacters map[PlayerCharacterType]PlayerCharacter

func loadPlayerCharacters() {
	PlayerCharacters = map[PlayerCharacterType]PlayerCharacter{
		PlayerCharacterKing: {
			Name:        "King",
			AddressedBy: "My Lord",
			Description: "The King is the ruler of the Kingdom.\n\nStarts with settlement level 2.",
			Image:       assets.PlayerKing,
		},
		PlayerCharacterQueen: {
			Name:        "Queen",
			AddressedBy: "My Queen",
			Description: "People love the Queen.\n\nStarts with +10 sight range.",
			Image:       assets.PlayerQueen,
		},
		PlayerCharacterPrincess: {
			Name:        "Princess",
			AddressedBy: "My Princess",
			Description: "The Princess will rule the kingdom.\n\nStart with 5 bonus food.",
			Image:       assets.PlayerPrincess,
		},
		PlayerCharacterMerchant: {
			Name:        "Merchant",
			AddressedBy: "Master Merchant",
			Description: "The Merchant is a master of trade.\n\nStarts with 10 bonus stone and wood.",
			Image:       assets.PlayerMerchant,
		},
	}
}
