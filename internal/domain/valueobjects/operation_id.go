package valueobjects

import (
	"github.com/google/uuid"
)

// OperationID representa o identificador único de uma operação
type OperationID struct {
	value string
}

// NewOperationID cria um novo ID de operação
func NewOperationID() OperationID {
	return OperationID{
		value: uuid.New().String(),
	}
}

// FromString cria um OperationID a partir de uma string
func OperationIDFromString(id string) (OperationID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return OperationID{}, err
	}
	return OperationID{value: id}, nil
}

// String retorna a representação em string
func (o OperationID) String() string {
	return o.value
}

// Equals verifica igualdade entre IDs
func (o OperationID) Equals(other OperationID) bool {
	return o.value == other.value
}
