package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/cszczepaniak/go-cribbly/internal/cribblyerr"
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

func (h *RequestHandler) HandleGetGameResult(ctx *gin.Context) {
	id := ctx.Param(`id`)
	r, err := h.pcfg.GameResultStore.Get(id)
	if cribblyerr.IsNotFound(err) {
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

	gameID, ok := ctx.Params.Get(`id`)
	if !ok {
		ctx.String(http.StatusBadRequest, `must specify associated game`)
		return
	}
	r.GameID = gameID

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

	var g model.Game
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		g, err = h.pcfg.GameStore.Get(r.GameID)
		return err
	})
	eg.Go(func() error {
		_, err := h.pcfg.TeamStore.Get(r.Winner)
		return err
	})

	err = eg.Wait()
	if cribblyerr.IsNotFound(err) {
		ctx.String(http.StatusBadRequest, fmt.Sprintf(`result cannot be made for non-existent resource: %v`, err))
		return
	} else if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	r.ID = r.GameID
	for _, t := range g.TeamIDs {
		if t == r.Winner {
			continue
		}
		r.Loser = t
		break
	}
	r, err = h.pcfg.GameResultStore.Create(r)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, r)
}
