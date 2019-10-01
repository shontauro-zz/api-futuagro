package store

import (
	"context"
	"log"
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

const supplierCollection string = "suppliers"

//MongoSupplierRepository a repo for saving suppliers into a mongo database
type MongoSupplierRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindByID returns a supplier by its ID from mongodb
func (repo *MongoSupplierRepository) FindByID(id string) (*models.Supplier, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var supplier *models.Supplier
	if err := result.Decode(&supplier); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a supplier")
	}
	return supplier, nil
}

// FindAll returns a list of suppliers from mongodb
func (repo *MongoSupplierRepository) FindAll() ([]*models.Supplier, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all suppliers")
	}

	var results []*models.Supplier
	for cursor.Next(context.TODO()) {
		var supplier models.Supplier
		if err := cursor.Decode(&supplier); err != nil {
			log.Printf("Error decoding a supplier on FindAll(): %v", err)
		} else {
			results = append(results, &supplier)
		}
	}
	err = cursor.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all suppliers")
	}
	return results, nil
}

// Insert a new supplier into mongodb
func (repo *MongoSupplierRepository) Insert(supplier *dtos.SupplierDto) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	now := primitive.DateTime(time.Now().UnixNano() / 1e6)
	data := bson.D{
		primitive.E{Key: "name", Value: supplier.Name},
		primitive.E{Key: "surname", Value: supplier.Surname},
		primitive.E{Key: "documentType", Value: supplier.DocumentType},
		primitive.E{Key: "documentNumber", Value: supplier.DocumentNumber},
		primitive.E{Key: "cityId", Value: supplier.CityID},
		primitive.E{Key: "email", Value: supplier.Email},
		primitive.E{Key: "addressLine1", Value: supplier.AddressLine1},
		primitive.E{Key: "phoneNumber", Value: supplier.PhoneNumber},
		primitive.E{Key: "createdAt", Value: now},
		primitive.E{Key: "updatedAt", Value: now},
		primitive.E{Key: "recordStatus", Value: enums.Active},
	}
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Inserting a new supplier")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update a supplier's document by its id in mongodb
func (repo *MongoSupplierRepository) Update(id string, supplier *dtos.SupplierDto) (*models.Supplier, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "name", Value: supplier.Name},
			primitive.E{Key: "surname", Value: supplier.Surname},
			primitive.E{Key: "documentType", Value: supplier.DocumentType},
			primitive.E{Key: "documentNumber", Value: supplier.DocumentNumber},
			primitive.E{Key: "cityId", Value: supplier.CityID},
			primitive.E{Key: "email", Value: supplier.Email},
			primitive.E{Key: "AddressLine1", Value: supplier.AddressLine1},
			primitive.E{Key: "phoneNumber", Value: supplier.PhoneNumber},
			primitive.E{Key: "updatedAt", Value: primitive.DateTime(time.Now().UnixNano() / 1e6)},
		},
	}}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedSupplier *models.Supplier
	if err := result.Decode(&updatedSupplier); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a supplier")
	}
	return updatedSupplier, nil
}

// Delete a supliers document from mongodb
func (repo *MongoSupplierRepository) Delete(id string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting a supplier")
	}
	return result.DeletedCount > 0, nil
}

//InsertCrop register a new crop for a supplier
func (repo *MongoSupplierRepository) InsertCrop(supplierID string, cropDto dtos.CropDto) (*models.Supplier, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	objID, err := primitive.ObjectIDFromHex(supplierID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	now := primitive.DateTime(time.Now().UnixNano() / 1e6)
	data := bson.D{
		primitive.E{Key: "_id", Value: primitive.NewObjectID()},
		primitive.E{Key: "cityId", Value: cropDto.CityID},
		primitive.E{Key: "plantingDate", Value: cropDto.PlantingDate},
		primitive.E{Key: "harvestDate", Value: cropDto.HarvestDate},
		primitive.E{Key: "variantId", Value: cropDto.VariantID},
		primitive.E{Key: "createdAt", Value: now},
		primitive.E{Key: "updatedAt", Value: now},
	}
	update := bson.D{
		primitive.E{
			Key: "$push",
			Value: bson.D{
				primitive.E{Key: "crops", Value: data},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedSupplier *models.Supplier
	if err := result.Decode(&updatedSupplier); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a supplier")
	}
	return updatedSupplier, nil
}

// NewMongoSupplierRepository returns a new instance of a MongoDB supplier repo.
func NewMongoSupplierRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoSupplierRepository {
	return &MongoSupplierRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
