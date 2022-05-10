package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cszczepaniak/go-cribbly/internal/cribblyerr"
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

func (h *RequestHandler) HandleGetGameResult(ctx *gin.Context) {
	id := ctx.Param(`id`)
	r, err := h.pcfg.GameResultStore.Get(id)
	if err == cribblyerr.ErrNotFound {
		ctx.String(http.StatusNotFound, `game result not found`)
		return
	} else if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, r)
}

func (h *RequestHandler) HandleGetAllGameResults(ctx *gin.Context) {
	rs, err := h.pcfg.GameResultStore.GetAll()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, rs)
}

func (h *RequestHandler) HandleCreateGameResult(ctx *gin.Context) {
	var r model.GameResult
	err := json.NewDecoder(ctx.Request.Body).Decode(&r)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if r.GameID == `` {
		ctx.String(http.StatusBadRequest, `must specify associated game`)
		return
	}

	if r.Winner == `` {
		ctx.String(http.StatusBadRequest, `must specify winning team ID`)
		return
	}

	if r.LoserScore < 0 || r.LoserScore > 120 {
		ctx.String(http.StatusBadRequest, `score difference must be between 0 and 120`)
		return
	}

	if r.ID != `` {
		ctx.String(http.StatusBadRequest, `cannot specify game result ID`)
		return
	}

	_, err = h.pcfg.GameStore.Get(r.GameID)
	if err == cribblyerr.ErrNotFound {
		ctx.String(http.StatusBadRequest, `result cannot be made for non-existent game`)
		return
	} else if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	r.ID = r.GameID
	r, err = h.pcfg.GameResultStore.Create(r)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, r)
}
