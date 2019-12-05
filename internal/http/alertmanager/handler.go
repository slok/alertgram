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
			ctx.AbortWithError(http.StatusBadRequest, err)
		}

		model, err := reqAlerts.toDomain()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}

		err = w.forwarder.Forward(ctx.Request.Context(), model)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
