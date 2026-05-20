package model

type RecordID string
type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string

type RatingValue int

type Rating struct {
	RecordID   string
	RecordType string
	UserID     UserID
	Value      RatingValue
}

// Type of rating event
type RatingEventType string

// An event containing rating information.
type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

// Rating event types
const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
