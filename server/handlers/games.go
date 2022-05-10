package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cszczepaniak/go-cribbly/internal/cribblyerr"
	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func (h *RequestHandler) HandleGetGame(ctx *gin.Context) {
	id := ctx.Param(`id`)
	g, err := h.pcfg.GameStore.Get(id)
	if err == cribblyerr.ErrNotFound {
		ctx.String(http.StatusNotFound, `game not found`)
		return
	} else if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, g)
}

func (h *RequestHandler) HandleGetAllGames(ctx *gin.Context) {
	gs, err := h.pcfg.GameStore.GetAll()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gs)
}

func (h *RequestHandler) HandleCreateGame(ctx *gin.Context) {
	var g model.Game
	err := json.NewDecoder(ctx.Request.Body).Decode(&g)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(g.TeamIDs) != 2 {
		ctx.String(http.StatusBadRequest, `expected two team IDs`)
		return
	}

	if g.ID != `` {
		ctx.String(http.StatusBadRequest, `cannot specify game ID`)
		return
	}

	if g.Kind != model.PrelimGame && g.Kind != model.TournamentGame {
		ctx.String(http.StatusBadRequest, `unknown game kind`)
	}

	g.ID = random.UUID()
	g, err = h.pcfg.GameStore.Create(g)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, g)
}
