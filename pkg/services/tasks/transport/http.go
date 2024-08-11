package transport

import (
	"net/http"

	db "github.com/maacarma/scheduler/pkg/db"
	svc "github.com/maacarma/scheduler/pkg/services/tasks"
	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	mongodb "github.com/maacarma/scheduler/pkg/services/tasks/store/mongodb"
	postgres "github.com/maacarma/scheduler/pkg/services/tasks/store/postgres"

	"github.com/gin-gonic/gin"
)

// Activate activates the router.
func Activate(router *gin.Engine, dbClients *db.Clients) {
	var repo svc.Repo
	switch {
	case dbClients.Pg != nil:
		repo = postgres.New(dbClients.Pg)
	case dbClients.Mongo != nil:
		repo = mongodb.New(dbClients.Mongo)
	}

	newHandler(router, svc.New(repo))
}

// handler is the concrete implementation of the tasks http methods.
type handler struct {
	service svc.Service
}

// newHandler creates a new handler
func newHandler(router *gin.Engine, sc svc.Service) {
	h := handler{
		service: sc,
	}
	router.GET("/tasks", h.GetAll)
	router.GET("/tasks/:namespace", h.GetAllByNamespace)
	router.POST("/tasks", h.CreateTask)
}

// GetAll returns all tasks
func (h *handler) GetAll(c *gin.Context) {
	tasks, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetAllByNamespace returns all tasks by namespace
func (h *handler) GetAllByNamespace(c *gin.Context) {
	tasks, err := h.service.GetByNamespace(c.Request.Context(), c.Param("namespace"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// CreateTask creates a new task
func (h *handler) CreateTask(c *gin.Context) {
	var task models.TaskPayload
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Create(c.Request.Context(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"id": id})
}
