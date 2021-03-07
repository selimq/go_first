package lib

import (
	"fmt"
	"time"
)

//Word ..
type Word struct {
	ID             int       `json:"ID" bson:"ID"`
	Text           string    `json:"Text" bson:"Text"`
	TranslatedText string    `json:"TranslatedText" bson:"TranslatedText"`
	CreatedAt      time.Time `json:"CreatedAt" bson:"CreatedAt"`
	Language       string    `json:"Language" bson:"Language"`
}

//Yazdir ..
func (w Word) Yazdir() {
	fmt.Println(w.Text + ": " + w.TranslatedText)
}

//Person is a common model
type Person struct {
	Ad string `json:"personName"`
	ID int    `json:"personID"`
}

//Customer ..
type Customer struct {
	Person
	Budget int
}
