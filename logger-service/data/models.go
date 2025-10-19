package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logger").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry:", err)
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logger").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error retrieving log entries:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var logEntry LogEntry
		err := cursor.Decode(&logEntry)
		if err != nil {
			log.Println("Error decoding log entry:", err)
			return nil, err
		}
		logs = append(logs, &logEntry)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}
	return logs, nil
}

func (l *LogEntry) GetByID(id string) (*LogEntry, error) {
	collection := client.Database("logger").Collection("logs")
	var logEntry LogEntry
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&logEntry)
	if err != nil {
		log.Println("Error retrieving log entry by ID:", err)
		return nil, err
	}
	return &logEntry, nil
}

func (l *LogEntry) DeleteByID(id string) error {
	collection := client.Database("logger").Collection("logs")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		log.Println("Error deleting log entry by ID:", err)
		return err
	}
	return nil
}

func (l *LogEntry) UpdateByID(id string, data string) error {
	collection := client.Database("logger").Collection("logs")
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"data", data},
				{"updated_at", time.Now()},
			}},
		},
	)
	if err != nil {
		log.Println("Error updating log entry by ID:", err)
		return err
	}
	return nil
}
func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logger").Collection("logs")
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection:", err)
		return err
	}
	return nil
}
