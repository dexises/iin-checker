package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Person entity for persistence.
type Person struct {
	IIN   string `bson:"iin,unique"`
	Name  string `bson:"name"`
	Phone string `bson:"phone"`
}

// PersonRepository defines persistence operations.
type PersonRepository interface {
	Create(ctx context.Context, p Person) error
	Exists(ctx context.Context, iin string) (bool, error)
	FindByName(ctx context.Context, namePart string) ([]Person, error)
	Get(ctx context.Context, iin string) (Person, error)
}

// mongoPersonRepo implements PersonRepository using MongoDB.
type mongoPersonRepo struct {
	col *mongo.Collection
}

// NewMongoPersonRepo constructs a PersonRepository backed by Mongo.
func NewMongoPersonRepo(db *mongo.Database) PersonRepository {
	col := db.Collection("persons")
	// ensure index on iin
	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "iin", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &mongoPersonRepo{col: col}
}

func (r *mongoPersonRepo) FindByName(ctx context.Context, namePart string) ([]Person, error) {
	filter := bson.M{"name": bson.M{"$regex": namePart, "$options": "i"}}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var results []Person
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	log.Printf("FindByName: found %d records for %q", len(results), namePart)

	return results, nil
}

// Create inserts a new Person.
func (r *mongoPersonRepo) Create(ctx context.Context, p Person) error {
	_, err := r.col.InsertOne(ctx, p)
	return err
}

// Exists checks if a Person with given IIN exists.
func (r *mongoPersonRepo) Exists(ctx context.Context, iin string) (bool, error) {
	filter := bson.M{"iin": iin}
	count, err := r.col.CountDocuments(ctx, filter)
	return count > 0, err
}

// List returns all Person records.
func (r *mongoPersonRepo) List(ctx context.Context) ([]Person, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var results []Person
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Get retrieves a Person by IIN.
func (r *mongoPersonRepo) Get(ctx context.Context, iin string) (Person, error) {
	var p Person
	filter := bson.M{"iin": iin}
	err := r.col.FindOne(ctx, filter).Decode(&p)
	return p, err
}
