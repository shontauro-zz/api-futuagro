package store

import (
	"context"
	"log"
	"time"

	"futuagro.com/pkg/config"
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const cropCollection string = "crops"

//MongoCropRepository a repo for saving crops into a mongo database
type MongoCropRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindByID returns a crop by its ID from mongodb
func (repo *MongoCropRepository) FindByID(id string) (*models.Crop, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cropCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	var pipeline = []bson.M{
		bson.M{"$match": bson.M{"_id": objID}},
	}
	pipeline = append(pipeline, buildStandardCropPipeline()...)

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	cursor, err := collection.Aggregate(ctx, pipeline, nil)

	var crop *models.Crop
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&crop); err != nil {
			log.Printf("Error decoding a crop: %v", err)
		}
		log.Printf("papu %v ", crop)
	}
	err = cursor.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a crop")
	}
	return crop, nil
}

// FindAll returns a list of crops from mongodb
func (repo *MongoCropRepository) FindAll() ([]*models.Crop, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cropCollection)
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	var pipeline = buildStandardCropPipeline()
	cursor, err := collection.Aggregate(ctx, pipeline, nil)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all crops")
	}

	var results []*models.Crop
	for cursor.Next(context.TODO()) {
		var crop models.Crop
		if err := cursor.Decode(&crop); err != nil {
			log.Printf("Error decoding a crop on FindAll(): %v", err)
		} else {
			results = append(results, &crop)
		}
	}

	err = cursor.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all crops")
	}
	return results, nil
}

// Insert a new crop into mongodb
func (repo *MongoCropRepository) Insert(dto *dtos.CropDto) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cropCollection)
	now := primitive.DateTime(time.Now().UnixNano() / 1e6)
	data := bson.D{
		primitive.E{Key: "cityId", Value: dto.CityID},
		primitive.E{Key: "plantingDate", Value: dto.PlantingDate},
		primitive.E{Key: "harvestDate", Value: dto.HarvestDate},
		primitive.E{Key: "variantId", Value: dto.VariantID},
		primitive.E{Key: "supplierId", Value: dto.SupplierID},
		primitive.E{Key: "createdAt", Value: now},
		primitive.E{Key: "updatedAt", Value: now},
	}
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Inserting a new crop")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update a crop document by its id in mongodb
func (repo *MongoCropRepository) Update(id string, dto *dtos.CropDto) (*models.Crop, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cropCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	update := bson.M{"$set": bson.M{
		"cityId":       dto.CityID,
		"plantingDate": dto.PlantingDate,
		"harvestDate":  dto.HarvestDate,
		"variantId":    dto.VariantID,
		"supplierId":   dto.SupplierID,
		"updatedAt":    time.Now(),
	}}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedCrop *models.Crop
	if err := result.Decode(&updatedCrop); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a crop")
	}
	return updatedCrop, nil
}

// Delete a supliers document from mongodb
func (repo *MongoCropRepository) Delete(id string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cropCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting a crop")
	}
	return result.DeletedCount > 0, nil
}

func buildStandardCropPipeline() []bson.M {
	return []bson.M{
		bson.M{"$lookup": bson.M{
			"from":         "cities",
			"localField":   "cityId",
			"foreignField": "_id",
			"as":           "city",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$city",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "variants",
			"localField":   "variantId",
			"foreignField": "_id",
			"as":           "variant",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$variant",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "items",
			"localField":   "variant.itemId",
			"foreignField": "_id",
			"as":           "variant.item",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$variant.item",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "suppliers",
			"localField":   "supplierId",
			"foreignField": "_id",
			"as":           "supplier",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$supplier",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$project": bson.M{
			"variantId":  0,
			"supplierId": 0,
			"cityId":     0,
		}},
	}
}

// NewMongoCropRepository returns a new instance of a MongoDB crop repo.
func NewMongoCropRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoCropRepository {
	return &MongoCropRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
