package alertmanager

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w webhookHandler) HandleAlerts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqAlerts := alertGroupV4{}
		err := ctx.BindJSON(&reqAlerts)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		}

		model, err := reqAlerts.toDomain()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		}

		err = w.forwarder.Forward(ctx.Request.Context(), model)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePublic)
		}
	}
}
