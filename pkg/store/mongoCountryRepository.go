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

const countryCollection = "countries"

// MongoCountryRepository a repository that implements the basic CRUD operations for saving countries into a mongo database
type MongoCountryRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindByID returns a country by its ID from mongodb
func (repo *MongoCountryRepository) FindByID(id string) (*models.Country, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var country *models.Country
	if err := result.Decode(&country); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error deconding a country")
	}
	return country, nil
}

// FindAll returns a list of countries from mongodb
func (repo *MongoCountryRepository) FindAll() ([]*models.Country, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	opts := options.Find().SetSort(bson.D{
		primitive.E{Key: "countryName", Value: 1},
		primitive.E{Key: "states.stateName", Value: 1},
	})
	cursor, err := collection.Find(context.Background(), bson.D{}, opts)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all countries")
	}

	var results []*models.Country = []*models.Country{}
	for cursor.Next(context.TODO()) {
		var country models.Country
		if err := cursor.Decode(&country); err != nil {
			log.Printf("Error decoding a country on findAll(): %v", err)
		} else {
			results = append(results, &country)
		}
	}
	err = cursor.Err()
	if err != nil {
		return results, errors.Wrap(err, "Error finding all countries")
	}
	return results, nil
}

// Insert a new country into mongodb
func (repo *MongoCountryRepository) Insert(country *models.Country) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	data := bson.D{
		primitive.E{Key: "countryName", Value: country.CountryName},
		primitive.E{Key: "countryCode", Value: country.CountryCode},
		primitive.E{Key: "states", Value: country.States},
		primitive.E{Key: "RecordStatus", Value: country.RecordStatus},
	}

	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Inserting a new country")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update a country by its id in mongodb
func (repo *MongoCountryRepository) Update(id string, country *models.Country) (*models.Country, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "countryName", Value: country.CountryName},
			primitive.E{Key: "countryCode", Value: country.CountryCode},
			primitive.E{Key: "RecordStatus", Value: country.RecordStatus},
		},
	}}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedCountry *models.Country
	if err := result.Decode(&updatedCountry); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a country")
	}
	return updatedCountry, nil
}

// Delete a country document from mongodb
func (repo *MongoCountryRepository) Delete(id string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting a country")
	}
	return result.DeletedCount > 0, nil
}

// InsertCountryState add a new state to a country
func (repo *MongoCountryRepository) InsertCountryState(countryID string, stateDto dtos.CountryStateDto) (*models.Country, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	objID, err := primitive.ObjectIDFromHex(countryID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	active := enums.Active
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	data := bson.D{
		primitive.E{Key: "countryState", Value: stateDto.StateName},
		primitive.E{Key: "RecordStatus", Value: &active},
	}
	update := bson.D{
		primitive.E{
			Key: "$push",
			Value: bson.D{
				primitive.E{Key: "states", Value: data},
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
	var updatedCountry *models.Country
	if err := result.Decode(&updatedCountry); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a country")
	}
	return updatedCountry, nil
}

// UpdateCountryState update the data of a country state
func (repo *MongoCountryRepository) UpdateCountryState(countryID string, stateID string, stateDto dtos.CountryStateDto) (*models.Country, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	countryObjID, err := primitive.ObjectIDFromHex(countryID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	stateObjID, err := primitive.ObjectIDFromHex(stateID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: countryObjID},
		primitive.E{Key: "states._id", Value: stateObjID},
	}
	data := bson.D{
		primitive.E{Key: "stateName", Value: stateDto.StateName},
		primitive.E{Key: "recordStatus", Value: stateDto.RecordStatus},
	}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "states.$", Value: data},
		},
	}}
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedCountry *models.Country
	if err := result.Decode(&updatedCountry); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a country")
	}
	return updatedCountry, nil
}

// DeleteCountryState remove a state from a country
func (repo *MongoCountryRepository) DeleteCountryState(countryID string, stateID string) (*models.Country, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	countryObjID, err := primitive.ObjectIDFromHex(countryID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	stateObjID, err := primitive.ObjectIDFromHex(stateID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{
		primitive.E{Key: "_id", Value: countryObjID},
	}
	stateData := bson.D{
		primitive.E{Key: "_id", Value: stateObjID},
	}
	update := bson.D{primitive.E{
		Key: "$pull",
		Value: bson.D{
			primitive.E{Key: "states", Value: stateData},
		},
	}}
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedCountry *models.Country
	if err := result.Decode(&updatedCountry); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a country")
	}
	return updatedCountry, nil
}

// NewMongoCountryRepository returns a new instance of a MongoDB country repository.
func NewMongoCountryRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoCountryRepository {
	return &MongoCountryRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
