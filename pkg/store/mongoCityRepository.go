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

const cityCollection = "cities"

// MongoCityRepository a repository that implements the basic CRUD operations for saving cities into a mongo database
type MongoCityRepository struct {
	databaseName string
	client       *mongo.Client
}

// FindByID returns a city by its ID from mongodb
func (repo *MongoCityRepository) FindByID(id string) (*models.City, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cityCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var city *models.City
	if err := result.Decode(&city); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error deconding a city")
	}
	return city, nil
}

// FindAll returns a list of cities from mongodb
func (repo *MongoCityRepository) FindAll() ([]*models.City, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cityCollection)
	cursor, err := collection.Find(context.Background(), bson.D{})
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "Error finding all cities")
	}
	cities, err := parseListOfCityDocs(cursor)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

//FindCitiesByCountryState find a list of cities by a country state ID
func (repo *MongoCityRepository) FindCitiesByCountryState(stateID string) ([]*models.City, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cityCollection)
	objID, err := primitive.ObjectIDFromHex(stateID)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "countryState", Value: objID}}
	cursor, err := collection.Find(context.Background(), filter)
	cities, err := parseListOfCityDocs(cursor)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

// Insert a new city into mongodb
func (repo *MongoCityRepository) Insert(dto *dtos.CityDto) (string, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cityCollection)
	stateID, err := primitive.ObjectIDFromHex(dto.CountryStateID)
	if err != nil {
		return string(""), errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	active := enums.Active
	data := bson.D{
		primitive.E{Key: "cityName", Value: dto.CityName},
		primitive.E{Key: "countryState", Value: stateID},
		primitive.E{Key: "recordStatus", Value: &active},
	}

	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return string(""), errors.Wrap(err, "Inserting a new city")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update a city's document by its id in mongodb
func (repo *MongoCityRepository) Update(id string, dto *dtos.CityDto) (*models.City, error) {
	collection := repo.client.Database(repo.databaseName).Collection(countryCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "cityName", Value: dto.CityName},
			primitive.E{Key: "countryState", Value: dto.CountryStateID},
			primitive.E{Key: "RecordStatus", Value: dto.RecordStatus},
		},
	}}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, updateOpts)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var updatedCity *models.City
	if err := result.Decode(&updatedCity); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error decoding a city")
	}
	return updatedCity, nil
}

// Delete a city document from mongodb
func (repo *MongoCityRepository) Delete(id string) (bool, error) {
	collection := repo.client.Database(repo.databaseName).Collection(cityCollection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing ObjectID from Hex")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, errors.Wrap(err, "Error deleting a city")
	}
	return result.DeletedCount > 0, nil
}

func parseListOfCityDocs(cursor *mongo.Cursor) ([]*models.City, error) {
	var results []*models.City
	for cursor.Next(context.TODO()) {
		var city models.City
		if err := cursor.Decode(&city); err != nil {
			log.Printf("Error decoding a city: %v", err)
		} else {
			results = append(results, &city)
		}
	}
	err := cursor.Err()
	if err != nil {
		return results, errors.Wrap(err, "Error parsing a list of cities")
	}
	return results, nil
}

// NewMongoCityRepository returns a new instance of a MongoDB country repository.
func NewMongoCityRepository(confPtr *config.Config, clientPtr *mongo.Client) *MongoCityRepository {
	return &MongoCityRepository{databaseName: confPtr.Database.Name, client: clientPtr}
}
