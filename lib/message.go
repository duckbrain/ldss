package lib

import (
	"fmt"
)

type Message interface {
	fmt.Stringer
	Finished() bool
}

type MessageDownload struct {
	err NotDownloadedErr
}

func (m MessageDownload) String() string {
	return fmt.Sprintf("Downloading \"%v\"", m.err)
}
func (m MessageDownload) Finished() bool {
	return false
}

type MessageDone struct {
	item Item
}

func (m MessageDone) String() string {
	return fmt.Sprintf("Finished loading \"%v\"", m.item)
}
func (m MessageDone) Item() Item {
	return m.item
}
func (m MessageDone) Finished() bool {
	return true
}

type MessageError struct {
	err error
}

func (m MessageError) String() string {
	return m.Error()
}
func (m MessageError) Error() string {
	return m.err.Error()
}
func (m MessageError) Finished() bool {
	return true
}
