package save_decoder

import "encoding/json"

type IntValue struct {
	Value int `json:"value"`
}

type IntArrayValue struct {
	Value []int `json:"value"`
}

type LethalSaveInfo struct {
	GroupCredits         IntValue      `json:"GroupCredits"`         // Argent disponible dans le store
	Stats_DaysSpent      IntValue      `json:"Stats_DaysSpent"`      // Nombre de jours total passés dans le jeu
	QuotaFulfilled       IntValue      `json:"QuotaFulfilled"`       // Quota déjà rempli
	ProfitQuota          IntValue      `json:"ProfitQuota"`          // Quota à atteindre
	CurrentPlanetID      IntValue      `json:"CurrentPlanetID"`      // ID de la planète actuelle
	ShipGrabbableItemIDs IntArrayValue `json:"shipGrabbableItemIDs"` // Liste des items au sol (scrap + store items)
	ShipScrapValues      IntArrayValue `json:"shipScrapValues"`      // Prix du scrap au sol
	DeadlineTime         IntValue      `json:"DeadlineTime"`         // Jours restants
}

func Read(jsonSave JsonSave) (LethalSaveInfo, error) {
	var saveInfo LethalSaveInfo
	err := json.Unmarshal(jsonSave.Data, &saveInfo)
	if err != nil {
		return LethalSaveInfo{}, err
	}

	return saveInfo, nil
}
