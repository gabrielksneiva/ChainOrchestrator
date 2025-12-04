package valueobjects

import "fmt"

// ChainType representa o tipo de blockchain
type ChainType string

const (
	ChainTypeEVM  ChainType = "EVM"
	ChainTypeBTC  ChainType = "BTC"
	ChainTypeTRON ChainType = "TRON"
	ChainTypeSOL  ChainType = "SOL"
)

// NewChainType cria e valida um novo ChainType
func NewChainType(value string) (ChainType, error) {
	ct := ChainType(value)
	if !ct.IsValid() {
		return "", fmt.Errorf("invalid chain type: %s", value)
	}
	return ct, nil
}

// IsValid verifica se o chain type é válido
func (c ChainType) IsValid() bool {
	switch c {
	case ChainTypeEVM, ChainTypeBTC, ChainTypeTRON, ChainTypeSOL:
		return true
	default:
		return false
	}
}

// String retorna a representação em string
func (c ChainType) String() string {
	return string(c)
}

// Equals verifica igualdade
func (c ChainType) Equals(other ChainType) bool {
	return c == other
}

// SNSTopicAttribute retorna o atributo para filtro SNS
func (c ChainType) SNSTopicAttribute() string {
	return string(c)
}
