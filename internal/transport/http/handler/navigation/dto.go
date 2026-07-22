package http_handler_navigation

import http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"

type NavigationResponseDTO struct {
	Cooldown http_dto.CooldownDTO `json:"cooldown"`
}
