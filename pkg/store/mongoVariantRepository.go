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

const variantCollection = "variants"

// MongoVariantRepository a repository that implements the basic CRUD operations for saving variants into a mongo database
type MongoVariantRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindVariantByID returns a Variant by its ID from mongodb
func (repo *MongoVariantRepository) FindVariantByID(ID string) (*models.Variant, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: objID},
	}
	return repo.findOneVariant(filter)
}

// FindOneVariantByItemID returns a Variant by its ID and Item ID from mongodb
func (repo *MongoVariantRepository) FindOneVariantByItemID(itemID string, variantID string) (*models.Variant, error) {
	objItemdID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	objVariantID, err := primitive.ObjectIDFromHex(variantID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: objVariantID},
		primitive.E{Key: "itemId", Value: objItemdID},
	}
	return repo.findOneVariant(filter)
}

func (repo *MongoVariantRepository) findOneVariant(filter interface{}) (*models.Variant, error) {
	collection := repo.client.Database(repo.databaseName).Collection(variantCollection)
	result := collection.FindOne(context.TODO(), filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var variant *models.Variant
	if err := result.Decode(&variant); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a Variant")
	}
	return variant, nil
}

// FindVariantsByItemID return a list of variants that belongs to a product from mongodb
func (repo *MongoVariantRepository) FindVariantsByItemID(itemID string) ([]*models.Variant, error) {
	collection := repo.client.Database(repo.databaseName).Collection(variantCollection)
	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{primitive.E{Key: "itemId", Value: objID}}
	cursor, err := collection.Find(context.Background(), filter)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all variants")
	}
	var results []*models.Variant = []*models.Variant{}
	for cursor.Next(context.TODO()) {
		var variant models.Variant
		if err := cursor.Decode(&variant); err != nil {
			log.Printf("Error decoding a Variant: %v", err)
		} else {
			results = append(results, &variant)
		}
	}

	err = cursor.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error finding variants by itemId "+itemID)
	}
	return results, nil
}

// Insert a new variant into mongodb
func (repo *MongoVariantRepository) Insert(itemID string, variantDto *dtos.VariantDto) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(variantCollection)
	objItemID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return string(""), errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	createdAt := primitive.DateTime(time.Now().UnixNano() / 1e6)
	active := enums.Active
	data := bson.D{
		primitive.E{Key: "name", Value: variantDto.Name},
		primitive.E{Key: "lname", Value: strings.ToLower(variantDto.Name)},
		primitive.E{Key: "itemId", Value: objItemID},
		primitive.E{Key: "createdAt", Value: createdAt},
		primitive.E{Key: "updatedAt", Value: createdAt},
		primitive.E{Key: "recordStatus", Value: active},
	}

	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Error inserting a new variant")
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update a variant's data by its id in mongodb
func (repo *MongoVariantRepository) Update(itemID string, variantID string, variantDto *dtos.VariantDto) (*models.Variant, error) {
	collection := repo.client.Database(repo.databaseName).Collection(variantCollection)
	objItemdID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	objVariantID, err := primitive.ObjectIDFromHex(variantID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: objVariantID},
		primitive.E{Key: "itemId", Value: objItemdID},
	}
	data := bson.D{
		primitive.E{Key: "name", Value: variantDto.Name},
		primitive.E{Key: "lname", Value: strings.ToLower(variantDto.Name)},
		primitive.E{Key: "updatedAt", Value: primitive.DateTime(time.Now().UnixNano() / 1e6)},
	}
	if variantDto.RecordStatus != nil {
		data = append(data, primitive.E{Key: "recordStatus", Value: variantDto.RecordStatus})
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
		return nil, errors.Wrap(result.Err(), "Error updating a Variant")
	}
	var updatedVariant *models.Variant
	if err := result.Decode(&updatedVariant); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a Variant")
	}
	return updatedVariant, nil
}

// Delete a variant document from mongodb
func (repo *MongoVariantRepository) Delete(itemID string, variantID string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(variantCollection)
	objItemID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	objVariantID, err := primitive.ObjectIDFromHex(variantID)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := primitive.D{
		primitive.E{Key: "_id", Value: objVariantID},
		primitive.E{Key: "itemId", Value: objItemID},
	}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting a Variant")
	}
	return result.DeletedCount > 0, nil
}

// NewMongoVariantRepository returns a new instance of a mongodb repository for variants.
func NewMongoVariantRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoVariantRepository {
	return &MongoVariantRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
