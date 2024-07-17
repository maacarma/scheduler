package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/maacarma/scheduler/pkg/services/scheduler"
	"github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres/sqlgen"
)

func Activate(router *gin.Engine, pgConn *pgx.Conn) {
	newHandler(router, scheduler.New(sqlgen.New(pgConn)))
}

type handler struct {
	service scheduler.Service
}

func newHandler(router *gin.Engine, sc scheduler.Service) {
	h := handler{
		service: sc,
	}
	router.GET("/tasks", h.GetAll)
}

func (h *handler) GetAll(c *gin.Context) {
	tasks, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
