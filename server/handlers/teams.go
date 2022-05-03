package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cszczepaniak/go-cribbly/internal/cribblyerr"
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

func (h *RequestHandler) HandleGetTeam(ctx *gin.Context) {
	id := ctx.Param(`id`)
	t, err := h.pcfg.TeamStore.Get(id)
	if err == cribblyerr.ErrNotFound {
		ctx.String(http.StatusNotFound, `team not found`)
		return
	} else if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, t)
}

func (h *RequestHandler) HandleGetAllTeams(ctx *gin.Context) {
	ts, err := h.pcfg.TeamStore.GetAll()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ts)
}

func (h *RequestHandler) HandleCreateTeam(ctx *gin.Context) {
	var t model.Team
	err := json.NewDecoder(ctx.Request.Body).Decode(&t)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(t.Players) != 2 {
		ctx.String(http.StatusBadRequest, `expected two players`)
		return
	}

	if t.ID != `` {
		ctx.String(http.StatusBadRequest, `cannot specify team ID`)
		return
	}

	t, err = h.pcfg.TeamStore.Create(t.Players[0], t.Players[1])
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, t)
}
