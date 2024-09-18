package save_decoder

type LethalSaveInfo struct {
	GroupCredits    int // Nombre d'argent disponible au shop
	Stats_DaysSpent int // Nombre de jours total passés dans le jeu
	ProfitQuota     int // Quota à atteindre
	QuotaFulfilled  int // Quota déjà rempli
	CurrentPlanetID int // ID de la planète actuelle
}

func FindData(rawJson []byte) {}
