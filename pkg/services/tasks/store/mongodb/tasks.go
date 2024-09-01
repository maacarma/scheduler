package mongodb

import (
	"context"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	utils "github.com/maacarma/scheduler/utils"

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

// GetByID returns a task from the database with the given id.
func (r *repo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	collection := r.client.Database(r.db).Collection(r.col)
	task := &models.Task{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetActiveTasks returns a list of active tasks
func (r *repo) GetActiveTasks(ctx context.Context, curUnix utils.Unix) ([]*models.Task, error) {
	collection := r.client.Database(r.db).Collection(r.col)
	cursor, err := collection.Find(ctx, bson.M{"paused": false, "end_unix": bson.M{"$gte": curUnix}})
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

// UpdateStatus updates the status of a task.
func (r *repo) UpdateStatus(ctx context.Context, id string, paused bool) error {
	collection := r.client.Database(r.db).Collection(r.col)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"paused": paused}})
	return err
}

// Delete deletes a task
func (r *repo) Delete(ctx context.Context, id string) error {
	collection := r.client.Database(r.db).Collection(r.col)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
