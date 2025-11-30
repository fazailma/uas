package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoAchievement represents achievement document stored in MongoDB
type MongoAchievement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID   string             `bson:"student_id" json:"student_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Date        string             `bson:"date" json:"date"` // Format: YYYY-MM-DD
	ProofURL    string             `bson:"proof_url" json:"proof_url"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at"`
}
