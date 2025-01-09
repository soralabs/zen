package id

import "github.com/google/uuid"

const (
	namespace = "zen"
)

var namespaceUUID = uuid.NewSHA1(uuid.NameSpaceOID, []byte(namespace))

type ID string

// ManagerID is a unique identifier for manager instances
type ManagerID string

func New() ID {
	id := uuid.New()
	return ID(id.String())
}

func FromString(s string) ID {
	id := uuid.NewSHA1(namespaceUUID, []byte(s))
	return ID(id.String())
}

func (id ID) String() string {
	return string(id)
}
