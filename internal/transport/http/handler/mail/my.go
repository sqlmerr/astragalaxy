package http_handler_mail

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetClaimablePackagesResponse struct {
	Data []PackageDTO `json:"data"`
}

func (h *MailHTTPHandler) GetClaimablePackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	agentID := auth.GetAgentIDFromContext(ctx)

	pkgs, err := h.service.GetClaimablePackages(ctx, agentID)
	if err != nil {
		rh.ErrorResponse(err, "Failed to get claimable mail packages")
		return
	}

	res := GetClaimablePackagesResponse{
		Data: lo.Map(pkgs, func(p model.MailPackage, _ int) PackageDTO {
			return convertModel(p)
		}),
	}
	rh.JSONResponse(http.StatusOK, res)
}
