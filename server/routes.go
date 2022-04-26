package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) registerRoutes() {
	s.eng.GET(`/ping`, func(ctx *gin.Context) {
		ctx.String(http.StatusOK, `pong`)
	})

	s.eng.POST(`/testdata/:key`, func(ctx *gin.Context) {
		k := ctx.Params.ByName(`key`)

		err := s.s3Client.Put(k, ctx.Request.Body)
		defer ctx.Request.Body.Close()
		if err != nil {
			ctx.String(http.StatusInternalServerError, `%v`, err)
			return
		}

		ctx.Status(http.StatusCreated)
	})

	s.eng.GET(`/testdata/:key`, func(ctx *gin.Context) {
		p := ctx.Params.ByName(`key`)

		resp, err := s.s3Client.Get(p)
		if err != nil {
			ctx.String(http.StatusInternalServerError, `%v`, err)
			return
		}

		ctx.String(http.StatusOK, `%s`, resp)
	})

	s.eng.GET(`/testdata/many/:prefix`, func(ctx *gin.Context) {
		p := ctx.Params.ByName(`prefix`)

		resp, err := s.s3Client.GetWithPrefix(p)
		if err != nil {
			ctx.String(http.StatusInternalServerError, `%v`, err)
			return
		}

		res := make(map[string]string, len(resp))
		for k, v := range resp {
			res[k] = string(v)
		}

		ctx.JSON(http.StatusOK, res)
	})
}
