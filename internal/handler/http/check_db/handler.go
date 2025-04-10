package metric

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc usecase.MetricUseCase
	l  logger.Interface
}

func New(uc usecase.MetricUseCase, l logger.Interface) *Handler {
	return &Handler{
		uc: uc,
		l:  l,
	}
}

func (h *Handler) Check(ctx *gin.Context) {
	err := h.uc.Check()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
