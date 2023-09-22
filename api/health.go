package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
 

func (server *Server) checkHeath(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
