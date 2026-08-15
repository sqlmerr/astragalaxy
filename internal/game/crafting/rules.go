package crafting_service

import (
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func CountTotalResourceVolume(resources []model.Resource) int {
	var amount int
	for _, r := range resources {
		amount += r.Amount
	}

	return amount
}
