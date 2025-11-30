package repository

import (
	"context"
	"errors"
	"time"

	"UAS/app/models"
	"UAS/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoAchievementRepository handles MongoDB achievement operations
type MongoAchievementRepository struct {
	collection *mongo.Collection
}

// NewMongoAchievementRepository creates a new instance
func NewMongoAchievementRepository() *MongoAchievementRepository {
	return &MongoAchievementRepository{
		collection: database.MongoDB.Collection("achievements"),
	}
}

// Create creates a new achievement document in MongoDB
func (r *MongoAchievementRepository) Create(ctx context.Context, achievement *models.MongoAchievement) (*models.MongoAchievement, error) {
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

// FindByID finds achievement by ID
func (r *MongoAchievementRepository) FindByID(ctx context.Context, id string) (*models.MongoAchievement, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid achievement id")
	}

	var achievement models.MongoAchievement
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("achievement not found")
		}
		return nil, err
	}

	return &achievement, nil
}

// FindByStudentID finds all achievements by student ID
func (r *MongoAchievementRepository) FindByStudentID(ctx context.Context, studentID string) ([]models.MongoAchievement, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"student_id": studentID,
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []models.MongoAchievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// Update updates an achievement document
func (r *MongoAchievementRepository) Update(ctx context.Context, id string, achievement *models.MongoAchievement) (*models.MongoAchievement, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid achievement id")
	}

	achievement.UpdatedAt = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": achievement},
	)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("achievement not found")
	}

	return achievement, nil
}

// SoftDelete soft deletes an achievement
func (r *MongoAchievementRepository) SoftDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid achievement id")
	}

	now := time.Now()
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"deleted_at": now}},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("achievement not found")
	}

	return nil
}

// FindAll finds all achievements (admin only, including deleted)
func (r *MongoAchievementRepository) FindAll(ctx context.Context) ([]models.MongoAchievement, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []models.MongoAchievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}
