package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *RequestHandler) HandleGetStandings(ctx *gin.Context) {
	standings, err := h.standingsService.GetStandings()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, standings)
}
