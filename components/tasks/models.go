package tasks

import "time"

type Task struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	Title          string    `json:"title" bson:"title,omitempty"`
	StartTime      time.Time `json:"startTime" bson:"startTime,omitempty"`
	EndTime        time.Time `json:"endTime" bson:"endTime,omitempty"`
	UserID         string    `json:"userId" bson:"userId,omitempty"`
	ReminderPeriod time.Time `json:"reminderPeriod" bson:"reminderPeriod,omitempty"`
	TimeAdded      time.Time `json:"-" bson:"timeAdded,omitempty"`
}
