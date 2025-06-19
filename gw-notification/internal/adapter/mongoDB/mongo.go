package mongoDB

import (
	"context"
	"fmt"
	"gw-notification/internal/domain"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository(db *mongo.Database, collectionName string) *EventRepository {
	return &EventRepository{
		collection: db.Collection(collectionName),
	}
}

func (m *EventRepository) SaveEvent(ctx context.Context, req domain.Event) error {

	_, err := m.collection.InsertOne(ctx, req)
	if err != nil {
		return fmt.Errorf("failed insert in mongo:%w", err)
	}

	return nil
}
