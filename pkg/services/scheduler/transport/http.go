package transport

import (
	"net/http"

	svc "github.com/maacarma/scheduler/pkg/services/scheduler"
	sqlgen "github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres/sqlgen"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func Activate(router *gin.Engine, pgConn *pgx.Conn) {
	newHandler(router, svc.New(sqlgen.New(pgConn)))
}

type handler struct {
	service svc.Service
}

func newHandler(router *gin.Engine, sc svc.Service) {
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
