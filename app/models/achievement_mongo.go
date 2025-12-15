package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoAchievement represents achievement document stored in MongoDB
// Collection: achievements (3.2.1)
// Document structure with flexible details field based on achievement type
type MongoAchievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	StudentID       string             `bson:"student_id" json:"student_id"`             // UUID reference to PostgreSQL
	AchievementType string             `bson:"achievement_type" json:"achievement_type"` // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`

	// Dynamic details field based on achievement type
	// Supports: competition, publication, organization, certification details + common fields
	Details map[string]interface{} `bson:"details" json:"details"`

	// File attachments with metadata
	Attachments []Attachment `bson:"attachments" json:"attachments"`

	// Tags for categorization and search
	Tags []string `bson:"tags" json:"tags"`

	// Points for achievement scoring
	Points int `bson:"points" json:"points"`

	// Timestamps
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // Soft delete
}

// Attachment represents file attachment metadata in achievements
type Attachment struct {
	FileName   string    `bson:"file_name" json:"file_name"`
	FileURL    string    `bson:"file_url" json:"file_url"`
	FileType   string    `bson:"file_type" json:"file_type"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}

// CompetitionDetails represents details for competition achievement
type CompetitionDetails struct {
	CompetitionName  string `bson:"competition_name" json:"competition_name"`
	CompetitionLevel string `bson:"competition_level" json:"competition_level"` // 'international', 'national', 'regional', 'local'
	Rank             int    `bson:"rank" json:"rank"`
	MedalType        string `bson:"medal_type" json:"medal_type"` // 'gold', 'silver', 'bronze'
	EventDate        string `bson:"event_date" json:"event_date"`
	Location         string `bson:"location" json:"location"`
	Organizer        string `bson:"organizer" json:"organizer"`
}

// PublicationDetails represents details for publication achievement
type PublicationDetails struct {
	PublicationType  string   `bson:"publication_type" json:"publication_type"` // 'journal', 'conference', 'book'
	PublicationTitle string   `bson:"publication_title" json:"publication_title"`
	Authors          []string `bson:"authors" json:"authors"`
	Publisher        string   `bson:"publisher" json:"publisher"`
	ISSN             string   `bson:"issn" json:"issn"`
	EventDate        string   `bson:"event_date" json:"event_date"`
	Score            int      `bson:"score" json:"score"`
}

// OrganizationDetails represents details for organization achievement
type OrganizationDetails struct {
	OrganizationName string `bson:"organization_name" json:"organization_name"`
	Position         string `bson:"position" json:"position"`
	StartDate        string `bson:"start_date" json:"start_date"`
	EndDate          string `bson:"end_date" json:"end_date"`
	Organizer        string `bson:"organizer" json:"organizer"`
	Location         string `bson:"location" json:"location"`
}

// CertificationDetails represents details for certification achievement
type CertificationDetails struct {
	CertificationName   string `bson:"certification_name" json:"certification_name"`
	IssuedBy            string `bson:"issued_by" json:"issued_by"`
	CertificationNumber string `bson:"certification_number" json:"certification_number"`
	ValidUntil          string `bson:"valid_until" json:"valid_until"`
	EventDate           string `bson:"event_date" json:"event_date"`
	Score               int    `bson:"score" json:"score"`
}
