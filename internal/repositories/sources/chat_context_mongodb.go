package sources

import (
	"chat_server_v2/internal/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoDBChatContextSource interacts with a MongoDB instance
type MongoDBChatContextSource struct {
	collection *mongo.Collection
}

func (r *MongoDBChatContextSource) CreateChatContext(chatContextObject *models.ChatContextObj) (*models.ChatContextObj, error) {
	// Insert the chat object into the MongoDB collection
	docID, err := r.collection.InsertOne(context.TODO(), chatContextObject)
	if err != nil {
		// Return an error if the insertion fails
		return nil, err
	}
	fmt.Printf("Inserted document with ID: %v\n", docID)
	// Return a success indicator (1 for success)
	return chatContextObject, nil
}

func (r *MongoDBChatContextSource) GetChatContext(chatId string) (*models.ChatContextObj, error) {
	var chatObject models.ChatContextObj
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": chatId}).Decode(&chatObject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Home session not found
		}
		return nil, err
	}
	return &chatObject, nil
}

func (r *MongoDBChatContextSource) ListChatContexts() ([]models.ChatContextObj, error) {
	chatContexts := make([]models.ChatContextObj, 0)

	// Find all documents in the chat collection
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a ChatContextObj
	for cursor.Next(context.TODO()) {
		var chatContex models.ChatContextObj
		if err := cursor.Decode(&chatContex); err != nil {
			return nil, err
		}
		chatContexts = append(chatContexts, chatContex)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return chatContexts, nil
}

func (r *MongoDBChatContextSource) UpdateChatContextEvents(chatId string, events []models.Event) (string, error) {
	filter := bson.M{"_id": chatId} // Filter based on chat ID
	update := bson.M{
		"$set": bson.M{
			"events": events,
			"count":  len(events),
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return chatId, err
}

func (r *MongoDBChatContextSource) SaveChat(chat *models.ChatContextObj) (string, error) {
	// Insert the chat object into the MongoDB collection
	_, err := r.collection.InsertOne(context.TODO(), chat)
	if err != nil {
		// Return an error if the insertion fails
		return "", err
	}

	// Return a success indicator (1 for success)
	return chat.ID, nil
}

func (r *MongoDBChatContextSource) DeleteChat(chatId string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// NewMongoDBChatContextSource creates a new instance of MongoDBChatContextSource
func NewMongoDBChatContextSource(uriString string, dbName string, collectionName string) *MongoDBChatContextSource {
	// Initialize MongoDB connection
	// mongodb+srv://evist-robot:PoMAk9FlsHxmEGbb@evistgcpdatabase.dxfk1.mongodb.net/?retryWrites=true&w=majority&appName=EvistGcpDatabase
	clientOptions := options.Client().ApplyURI(uriString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection by pinging the server
	/*
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}
		log.Println("Connected to MongoDB!")


	*/
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDBChatContextSource{collection: collection}
}
