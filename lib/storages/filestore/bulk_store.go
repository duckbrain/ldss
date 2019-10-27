package filestore

// bulkStore implements the Store interface, but does all of its operations in a single transaction
type bulkStore struct {
}
