package sources

import (
	"chat_server_v2/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoDBChatContextRepository interacts with a MongoDB instance
type MongoDBChatScenarioStateSource struct {
	collection *mongo.Collection
}

func (source *MongoDBChatScenarioStateSource) CreateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error) {
	// Insert the chat object into the MongoDB collection
	_, err := source.collection.InsertOne(context.TODO(), obj)
	if err != nil {
		// Return an error if the insertion fails
		return nil, err
	}

	// Return a success indicator (1 for success)
	return obj, nil
}

func (source *MongoDBChatScenarioStateSource) GetChatScenarioState(chatId string) (*models.ChatScenarioStateObj, error) {
	var dialogueStateObject models.ChatScenarioStateObj
	err := source.collection.FindOne(context.TODO(), bson.M{"_id": chatId}).Decode(&dialogueStateObject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Home session not found
		}
		return nil, err
	}
	return &dialogueStateObject, nil
}

func (source *MongoDBChatScenarioStateSource) ListChatScenarioStates() ([]models.ChatScenarioStateObj, error) {
	chatScenarioStates := make([]models.ChatScenarioStateObj, 0)
	// Find all documents in the chat collection
	cursor, err := source.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a ChatContextObj
	for cursor.Next(context.TODO()) {
		var dialogueStateObject models.ChatScenarioStateObj
		if err := cursor.Decode(&dialogueStateObject); err != nil {
			return nil, err
		}
		chatScenarioStates = append(chatScenarioStates, dialogueStateObject)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return chatScenarioStates, nil
}

func (source *MongoDBChatScenarioStateSource) UpdateChatScenarioState(obj *models.ChatScenarioStateObj) (*models.ChatScenarioStateObj, error) {
	filter := bson.M{"_id": obj.ChatId}                  // Filter based on chat ID
	update := bson.M{"$set": bson.M{"state": obj.State}} // Update the messages field

	_, err := source.collection.UpdateOne(context.TODO(), filter, update)
	return obj, err
}

func NewMongoDBChatStateSource(uriString string, dbName string, collectionName string) (*MongoDBChatScenarioStateSource, error) {
	// Initialize MongoDB connection
	clientOptions := options.Client().ApplyURI(uriString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDBChatScenarioStateSource{collection: collection}, nil
}
