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
	"golang.org/x/crypto/bcrypt"
)

const userCollection string = "users"

//MongoUserRepository a repo for saving users into a mongo database
type MongoUserRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindByID returns an user by its ID from mongodb
func (repo *MongoUserRepository) FindByID(id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	user, err := repo.findOneUserBy(filter)
	if user != nil {
		user.HashedPassword = ""
	}
	return user, nil
}

// FindByEmail returns an user by its ID from mongodb
func (repo *MongoUserRepository) FindByEmail(email string) (*models.User, error) {

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	user, err := repo.findOneUserBy(filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *MongoUserRepository) findOneUserBy(filter interface{}) (*models.User, error) {
	collection := repo.client.Database(repo.databaseName).Collection(userCollection)
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user *models.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding an user")
	}
	return user, nil
}

// PopulateUserByID return an user with the crops property populated with the variants data
func (repo *MongoUserRepository) PopulateUserByID(id string) (*models.User, error) {
	collection := repo.client.Database(repo.databaseName).Collection(userCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	var pipeline = []bson.M{
		bson.M{"$match": bson.M{"_id": objID}},
	}
	pipeline = append(pipeline, buildStandardUserPipeline()...)

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	cursor, err := collection.Aggregate(ctx, pipeline, nil)

	var user *models.User
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Error decoding an user: %v", err)
		}
	}
	err = cursor.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding an user in PopulateUserByID")
	}
	return user, nil
}

// FindAll returns a list of users from mongodb
func (repo *MongoUserRepository) FindAll() ([]*models.User, error) {
	collection := repo.client.Database(repo.databaseName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	var pipeline = buildStandardUserPipeline()
	cursor, err := collection.Aggregate(ctx, pipeline, nil)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all users")
	}

	var results []*models.User
	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Error decoding an user on FindAll(): %v", err)
		} else {
			results = append(results, &user)
		}
	}
	err = cursor.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all users")
	}
	return results, nil
}

// Insert a new user into mongodb
func (repo *MongoUserRepository) Insert(dto *dtos.UserDto) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(userCollection)
	now := primitive.DateTime(time.Now().UnixNano() / 1e6)

	hashedPwdInBytes, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return string(""), errors.Wrap(err, "hashing a password")
	}

	data := bson.D{
		primitive.E{Key: "name", Value: dto.Name},
		primitive.E{Key: "surname", Value: dto.Surname},
		primitive.E{Key: "documentType", Value: dto.DocumentType},
		primitive.E{Key: "documentNumber", Value: dto.DocumentNumber},
		primitive.E{Key: "cityId", Value: dto.CityID},
		primitive.E{Key: "email", Value: dto.Email},
		primitive.E{Key: "hashedPassword", Value: string(hashedPwdInBytes)},
		primitive.E{Key: "addressLine1", Value: dto.AddressLine1},
		primitive.E{Key: "phoneNumber", Value: dto.PhoneNumber},
		primitive.E{Key: "role", Value: "user"},
		primitive.E{Key: "createdAt", Value: now},
		primitive.E{Key: "updatedAt", Value: now},
		primitive.E{Key: "recordStatus", Value: enums.Active},
	}
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Inserting a new user")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update an user document by its id in mongodb
func (repo *MongoUserRepository) Update(id string, dto *dtos.UserDto) (*models.User, error) {
	collection := repo.client.Database(repo.databaseName).Collection(userCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "name", Value: dto.Name},
			primitive.E{Key: "surname", Value: dto.Surname},
			primitive.E{Key: "documentType", Value: dto.DocumentType},
			primitive.E{Key: "documentNumber", Value: dto.DocumentNumber},
			primitive.E{Key: "cityId", Value: dto.CityID},
			primitive.E{Key: "email", Value: dto.Email},
			primitive.E{Key: "addressLine1", Value: dto.AddressLine1},
			primitive.E{Key: "phoneNumber", Value: dto.PhoneNumber},
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
	var updatedUser *models.User
	if err := result.Decode(&updatedUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding an user")
	}
	return updatedUser, nil
}

// Delete a supliers document from mongodb
func (repo *MongoUserRepository) Delete(id string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(userCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting an user")
	}
	return result.DeletedCount > 0, nil
}

func buildStandardUserPipeline() []bson.M {
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
			"role": bson.M{
				"$first": "$role",
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
			"role":           1,
			"crops": bson.M{
				"$cond": bson.A{bson.M{"$eq": bson.A{"$hasCrops", 0}}, bson.A{}, "$crops"},
			},
		}},
	}
}

// NewMongoUserRepository returns a new instance of a MongoDB user repo.
func NewMongoUserRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoUserRepository {
	return &MongoUserRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
