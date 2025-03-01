package utils

var validPlanetThreats = map[string]bool{
	"RADIATION": true,
	"TOXINS":    true,
	"FREEZING":  true,
	"HEAT":      true,
}

func ValidatePlanetThreat(threat string) bool {
	return validPlanetThreats[threat]
}
