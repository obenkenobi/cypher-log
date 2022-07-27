package database

import "github.com/kamva/mgm/v3"

// MongoModel An interface representing functionalities MongoDB model structures must have
type MongoModel interface {
	mgm.Model
	mgm.CollectionNameGetter
	GetIdStr() string // Gets the ID as a string. If the ID is not set, an empty string will be returned
	IsIdEmpty() bool  // Checks if the ID is not set
}
