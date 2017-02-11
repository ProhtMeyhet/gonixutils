package hashsum

import(
	"hash"

	"crypto/sha1"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
)

type Type uint8

// is valid
func (t Type) Valid() bool {
	return t > 0 && t <= CRC64
}

// return corresponding hash factory
func (t Type) Factory(input *Input) func() hash.Hash {
	switch(input.Type) {
	case MD5:
		return md5.New
	case SHA1:
		return sha1.New
	case SHA224:
		return sha256.New224
	case SHA256:
		return sha256.New
	case SHA384:
		return sha512.New384
	case SHA512:
		return sha512.New
	case SHA512_224:
		return sha512.New512_224
	case SHA512_256:
		return sha512.New512_256
	case ADLER32:
		return func() hash.Hash {
			return adler32.New()
		}
	case CRC32:
		return func() hash.Hash {
			return crc32.NewIEEE()
		}
	case CRC64:
		return func() hash.Hash {
			return crc64.New(crc64.MakeTable(crc64.ISO))
		}
	}

	return nil
}

// string me!
func (t Type) String() string {
	switch(t) {
	case MD5:
		return "MD5"
	case SHA1:
		return "SHA1"
	case SHA224:
		return "SHA224"
	case SHA256:
		return "SHA256"
	case SHA384:
		return "SHA384"
	case SHA512:
		return "SHA512"
	case SHA512_224:
		return "SHA512_224"
	case SHA512_256:
		return "SHA512_256"
	case ADLER32:
		return "ADLER32"
	case CRC32:
		return "CRC32"
	case CRC64:
		return "CRC64"
	case NONE:
		return "NONE"
	}

	return "INVALID"
}

// return the length of the hex encoded hash in byte
func (t Type) Len() uint {
	switch(t) {
	case MD5:
		return 32
	case SHA1:
		return 40
	case SHA224:
		return 56
	case SHA256:
                return 64
        case SHA384:
                return 96
        case SHA512:
                return 128
        case SHA512_224:
                return 56
        case SHA512_256:
                return 64
        case ADLER32:
		return 8
        case CRC32:
		return 8
        case CRC64:
		return 16
	}

	return 0
}
