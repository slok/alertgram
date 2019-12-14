package alertmanager

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/internalerrors"
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

		props := forward.Properties{
			CustomChatID: ctx.Query(w.cfg.ChatIDQueryString),
		}
		err = w.forwarder.Forward(ctx.Request.Context(), props, model)
		if err != nil {
			w.logger.Errorf("error forwarding alert: %s", err)

			if errors.Is(err, internalerrors.ErrInvalidConfiguration) {
				_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
				return
			}

			_ = ctx.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePublic)
			return
		}
	}
}
