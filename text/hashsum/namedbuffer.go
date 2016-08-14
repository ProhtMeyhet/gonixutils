package hash

import(

)

type NamedBuffer struct {
	name string
	buffer []byte
	done bool
	read int
}

func NewNamedBuffer(aname string) (namedBuffer NamedBuffer) {
	namedBuffer = NamedBuffer{ name: aname }
	return
}
