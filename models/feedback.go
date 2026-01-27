package models

type Feedback Article

type Reaction struct {
	ID              int      `json:"id" gorm:"primaryKey"`
	Rating          string   `json:"rating"`
	Comment         string   `json:"comment"`
	FeedbackID      int      `json:"feedback_id"`
	FeedbackProject string   `json:"feedback_project"`
	Feedback        Feedback `gorm:"foreignKey:FeedbackID,FeedbackProject;references:ID,Project"`
}

type RateQuery struct {
	Article Article `json:"article"`
	Rating  string  `json:"rating"`
}

type CommentQuery struct {
	ID      int    `json:"id"`
	Project string `json:"project"`
	Comment string `json:"comment"`
}
