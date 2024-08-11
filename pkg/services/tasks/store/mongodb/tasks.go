package mongodb

import (
	"context"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	client *mongo.Client
	db     string
	col    string
}

// New returns a new instance of the postgres repo.
func New(client *mongo.Client) *repo {
	return &repo{client: client, db: "scheduler", col: "tasks"}
}

// GetAll returns all tasks from the database.
func (r *repo) GetAll(ctx context.Context) ([]*models.Task, error) {
	collection := r.client.Database(r.db).Collection(r.col)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tasks := []*models.Task{}
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByNamespace returns all tasks from the database with the given namespace.
func (r *repo) GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error) {
	collection := r.client.Database(r.db).Collection(r.col)
	cursor, err := collection.Find(ctx, bson.M{"namespace": namespace})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tasks := []*models.Task{}
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// CreateOne creates a new task and returns the id.
func (r *repo) CreateOne(ctx context.Context, task *models.TaskPayload) (string, error) {
	collection := r.client.Database(r.db).Collection(r.col)
	res, err := collection.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
