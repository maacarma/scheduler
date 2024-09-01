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
func Activate(router *gin.Engine, dbClients *db.Clients, scheduler svc.Scheduler) {
	var repo svc.Repo
	switch {
	case dbClients.Pg != nil:
		repo = postgres.New(dbClients.Pg)
	case dbClients.Mongo != nil:
		repo = mongodb.New(dbClients.Mongo)
	}

	newHandler(router, svc.New(repo, scheduler))
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
	router.POST("/tasks", h.CreateTask)
	router.DELETE("/tasks/:id", h.DeleteTask)
	router.PUT("/tasks/:id/status", h.ToggleStatus)
	router.GET("/tasks/n/:namespace", h.GetAllByNamespace)
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

	if err := task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, statusCode, err := h.service.Create(c.Request.Context(), &task)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"id": id})
}

// ToggleStatus toggle the status of a task
func (h *handler) ToggleStatus(c *gin.Context) {
	id := c.Param("id")
	err := h.service.ToggleStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]bool{"updated": true})
}

func (h *handler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]bool{"deleted": true})
}
