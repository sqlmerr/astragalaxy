package v1

import (
	"astragalaxy/internal/service"
	"astragalaxy/internal/state"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	s     *service.Service
	state *state.State
}

func NewHandler(state *state.State) Handler {
	return Handler{
		s:     state.S,
		state: state,
	}
}

func (h *Handler) Register(api huma.API) {
	auth := huma.NewGroup(api, "/auth")
	h.registerAuthGroup(auth)

	h.registerSpaceshipsGroup(huma.NewGroup(api, "/spaceships"))
	h.registerAstralsGroup(huma.NewGroup(api, "/astral"))
	h.registerSystemsGroup(huma.NewGroup(api, "/systems"))
	h.registerRegistryGroup(huma.NewGroup(api, "/registry"))
	h.registerPlanetsGroup(huma.NewGroup(api, "/planets"))
	h.registerNavigationGroup(huma.NewGroup(api, "/navigation"))
	h.registerInventoryGroup(huma.NewGroup(api, "/inventory"))
	h.registerWalletGroup(huma.NewGroup(api, "/wallets"))
	h.registerExplorationGroup(huma.NewGroup(api, "/explore"))

}
