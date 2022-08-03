package webservices

// Server is a general server interface
type Server interface {
	// Run initiates the server and listens now for underlying network requests.
	Run()
}
