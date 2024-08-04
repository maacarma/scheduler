package transport

import (
	"net/http"

	svc "github.com/maacarma/scheduler/pkg/services/scheduler"
	models "github.com/maacarma/scheduler/pkg/services/scheduler/models"
	postgres "github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func Activate(router *gin.Engine, pgConn *pgx.Conn) {
	var repo svc.Repo
	switch {
	case pgConn != nil:
		repo = postgres.New(pgConn)
	}
	newHandler(router, svc.New(repo))
}

type handler struct {
	service svc.Service
}

func newHandler(router *gin.Engine, sc svc.Service) {
	h := handler{
		service: sc,
	}
	router.GET("/tasks", h.GetAll)
	router.GET("/tasks/:namespace", h.GetAllByNamespace)
	router.POST("/tasks", h.CreateTask)
}

func (h *handler) GetAll(c *gin.Context) {
	tasks, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *handler) GetAllByNamespace(c *gin.Context) {
	tasks, err := h.service.GetByNamespace(c.Request.Context(), c.Param("namespace"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *handler) CreateTask(c *gin.Context) {
	var task models.TaskPayload
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateTask(c.Request.Context(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, map[string]int64{"id": id})
}
