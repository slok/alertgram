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
			w.logger.Errorf("error unmarshalling JSON: %s", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
			return
		}

		model, err := reqAlerts.toDomain()
		if err != nil {
			w.logger.Errorf("error mapping to domain models: %s", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
			return
		}

		err = w.forwarder.Forward(ctx.Request.Context(), model)
		if err != nil {
			w.logger.Errorf("error forwarding alert: %s", err)
			_ = ctx.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePublic)
			return
		}
	}
}
