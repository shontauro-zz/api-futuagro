package store

import (
	"context"
	"log"
	"strings"
	"time"

	"futuagro.com/pkg/config"
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/enums"
	"futuagro.com/pkg/domain/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const itemCollection = "items"

// MongoItemRepository a repository that implements the basic CRUD operations for saving items into a mongo database
type MongoItemRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindByID returns an Item by its ID from mongodb
func (repo *MongoItemRepository) FindByID(id string) (*models.Item, error) {
	collection := repo.client.Database(repo.databaseName).Collection(itemCollection)
	objdID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objdID}}
	result := collection.FindOne(context.TODO(), filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var item *models.Item
	if err := result.Decode(&item); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding an Item")
	}
	return item, nil
}

// FindAll return a list of items from mongodb
func (repo *MongoItemRepository) FindAll() ([]*models.Item, error) {
	collection := repo.client.Database(repo.databaseName).Collection(itemCollection)
	cursor, err := collection.Find(context.Background(), bson.D{})
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all items")
	}
	var results []*models.Item = []*models.Item{}
	for cursor.Next(context.TODO()) {
		var item models.Item
		if err := cursor.Decode(&item); err != nil {
			log.Printf("Error decoding an Item: %v", err)
		} else {
			results = append(results, &item)
		}
	}

	err = cursor.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all item")
	}
	return results, nil
}

// Insert a new Item into mongodb
func (repo *MongoItemRepository) Insert(itemDto *dtos.ItemDto) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(itemCollection)
	createdAt := primitive.DateTime(time.Now().UnixNano() / 1e6)
	active := enums.Active
	data := bson.D{
		primitive.E{Key: "name", Value: itemDto.Name},
		primitive.E{Key: "lname", Value: strings.ToLower(itemDto.Name)},
		primitive.E{Key: "createdAt", Value: createdAt},
		primitive.E{Key: "updatedAt", Value: createdAt},
		primitive.E{Key: "recordStatus", Value: active},
	}

	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Error inserting a new Item")
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update an item's data by its id in mongodb
func (repo *MongoItemRepository) Update(id string, itemDto *dtos.ItemDto) (*models.Item, error) {
	collection := repo.client.Database(repo.databaseName).Collection(itemCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	data := bson.D{
		primitive.E{Key: "name", Value: itemDto.Name},
		primitive.E{Key: "lname", Value: strings.ToLower(itemDto.Name)},
		primitive.E{Key: "updatedAt", Value: primitive.DateTime(time.Now().UnixNano() / 1e6)},
	}
	if itemDto.RecordStatus != nil {
		data = append(data, primitive.E{Key: "recordStatus", Value: itemDto.RecordStatus})
	}
	update := bson.D{primitive.E{
		Key:   "$set",
		Value: data,
	}}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, errors.Wrap(result.Err(), "Error updating an Item")
	}
	var updatedItem *models.Item
	if err := result.Decode(&updatedItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding an Item")
	}
	return updatedItem, nil
}

// Delete an item document from mongodb
func (repo *MongoItemRepository) Delete(id string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(itemCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := primitive.D{
		primitive.E{Key: "_id", Value: objID},
	}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting an item")
	}
	return result.DeletedCount > 0, nil
}

// NewMongoItemRepository returns a new instance of a mongodb repository for items.
func NewMongoItemRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoItemRepository {
	return &MongoItemRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
