package hashsum

import(
	"hash"
)

type TypeInterface interface {
	Valid() bool
	String() string
	Len() uint
	Factory() func() hash.Hash
}
