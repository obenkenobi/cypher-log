package datasource

import "github.com/kamva/mgm/v3"

// MongoModel An interface representing functionalities MongoDB model structures must have
type MongoModel interface {
	mgm.Model
	mgm.CollectionNameGetter
	// GetIdStr Gets the ID as a string. If the ID is not set, an empty string will
	// be returned
	GetIdStr() string
	// IsIdEmpty Checks if the ID is not set
	IsIdEmpty() bool
}
