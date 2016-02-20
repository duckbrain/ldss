package lib

import (
	"fmt"
)

// Represents a message that comes from attempting a lookup.
type Message interface {
	fmt.Stringer
}

// Message indicating that a download needs to occur before proceding
type MessageDownload struct {
	err NotDownloadedErr
}

// Use facing message, indicating the download
func (m MessageDownload) String() string {
	return fmt.Sprintf("Downloading \"%v\"", m.err.String())
}

// Message indciating a result has been found
type MessageDone struct {
	item interface{}
}

// User facing message, indicating what was found.
func (m MessageDone) String() string {
	return fmt.Sprintf("Finished loading \"%v\"", m.item)
}

// The found result
func (m MessageDone) Item() interface{} {
	return m.item
}

// Message indicating an error occured
type MessageError struct {
	err error
}

// Outputs the error
func (m MessageError) String() string {
	return m.Error()
}

// Outputs the error
func (m MessageError) Error() string {
	return m.err.Error()
}
