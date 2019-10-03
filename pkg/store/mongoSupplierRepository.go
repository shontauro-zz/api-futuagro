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

// PopulateSupplierByID return a supplier with the crops property populated with the variant data
func (repo *MongoSupplierRepository) PopulateSupplierByID(id string) (*models.Supplier, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	var pipeline = []bson.M{
		bson.M{"$match": bson.M{"_id": objID}},
	}
	pipeline = append(pipeline, buildStandardSupplierPipeline()...)

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	cursor, err := collection.Aggregate(ctx, pipeline, nil)

	var supplier *models.Supplier
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&supplier); err != nil {
			log.Printf("Error decoding a supplier: %v", err)
		}
	}
	err = cursor.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a supplier in PopulateSupplierByID")
	}
	return supplier, nil
}

// FindAll returns a list of suppliers from mongodb
func (repo *MongoSupplierRepository) FindAll() ([]*models.Supplier, error) {
	collection := repo.client.Database(repo.databaseName).Collection(supplierCollection)
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	var pipeline = buildStandardSupplierPipeline()
	cursor, err := collection.Aggregate(ctx, pipeline, nil)
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

func buildStandardSupplierPipeline() []bson.M {
	return []bson.M{
		bson.M{"$lookup": bson.M{
			"from":         "crops",
			"localField":   "_id",
			"foreignField": "supplierId",
			"as":           "crops",
		}},
		bson.M{"$addFields": bson.M{
			"hasCrops": bson.M{
				"$size": bson.M{
					"$ifNull": bson.A{"$crops", bson.A{}}}}}},
		bson.M{"$unwind": bson.M{
			"path":                       "$crops",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "variants",
			"localField":   "crops.variantId",
			"foreignField": "_id",
			"as":           "crops.variant",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$crops.variant",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "items",
			"localField":   "crops.variant.itemId",
			"foreignField": "_id",
			"as":           "crops.variant.item",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$crops.variant.item",
			"preserveNullAndEmptyArrays": true,
		}},
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
		bson.M{"$group": bson.M{
			"_id": "$_id",
			"name": bson.M{
				"$first": "$name",
			},
			"surname": bson.M{
				"$first": "$surname",
			},
			"documentType": bson.M{
				"$first": "$documentType",
			},
			"documentNumber": bson.M{
				"$first": "$documentNumber",
			},
			"city": bson.M{
				"$first": "$city",
			},
			"email": bson.M{
				"$first": "$email",
			},
			"addressLine1": bson.M{
				"$first": "$addressLine1",
			},
			"phoneNumber": bson.M{
				"$first": "$phoneNumber",
			},
			"createdAt": bson.M{
				"$first": "$createdAt",
			},
			"updatedAt": bson.M{
				"$first": "$updatedAt",
			},
			"recordStatus": bson.M{
				"$first": "$recordStatus",
			},
			"hasCrops": bson.M{
				"$first": "$hasCrops",
			},
			"crops": bson.M{
				"$push": "$crops",
			},
		}},
		bson.M{"$project": bson.M{
			"_id":            1,
			"name":           1,
			"surname":        1,
			"documentType":   1,
			"documentNumber": 1,
			"city":           1,
			"email":          1,
			"addressLine1":   1,
			"phoneNumber":    1,
			"createdAt":      1,
			"updatedAt":      1,
			"recordStatus":   1,
			"crops": bson.M{
				"$cond": bson.A{bson.M{"$eq": bson.A{"$hasCrops", 0}}, bson.A{}, "$crops"},
			},
		}},
	}
}

// NewMongoSupplierRepository returns a new instance of a MongoDB supplier repo.
func NewMongoSupplierRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoSupplierRepository {
	return &MongoSupplierRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
