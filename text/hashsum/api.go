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
//	"hash/fnv"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// hash a file or stdin
func Hash(input *Input) (exitCode uint8) {
	if input.Compare {
		return Compare(input)
	}

	switch(input.Type) {
	case MD5:
		return Md5(input)
	case SHA1:
		return Sha1(input)
	case SHA224:
		return Sha224(input)
	case SHA256:
		return Sha256(input)
	case SHA384:
		return Sha384(input)
	case SHA512:
		return Sha512(input)
	case SHA512_224:
		return Sha512_224(input)
	case SHA512_256:
		return Sha512_256(input)
	case ADLER32:
		return Adler32(input)
	case CRC32:
		return Crc32(input)
	case CRC64:
		return Crc64(input)
	}

	return abstract.ERROR_NO_INPUT
}

// return corresponding hash factory
func Factory(input *Input) func() hash.Hash {
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

// hash a file or stdin
func Md5(input *Input) (exitCode uint8) {
	return doHash(input, md5.New)
}

// hash a file or stdin
func Sha1(input *Input) (exitCode uint8) {
	return doHash(input, sha1.New)
}

// hash a file or stdin
func Sha224(input *Input) (exitCode uint8) {
	return doHash(input, sha256.New224)
}

// hash a file or stdin
func Sha256(input *Input) (exitCode uint8) {
	return doHash(input, sha256.New)
}

// hash a file or stdin
func Sha384(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New384)
}

// hash a file or stdin
func Sha512(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New)
}

// hash a file or stdin
func Sha512_224(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New512_224)
}
// hash a file or stdin
func Sha512_256(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New512_256)
}
// hash a file or stdin
func Adler32(input *Input) (exitCode uint8) {
	return doHash(input, func() hash.Hash {
		return adler32.New()
	})
}
// hash a file or stdin
func Crc32(input *Input) (exitCode uint8) {
	return doHash(input, func() hash.Hash {
		return crc32.NewIEEE()
	})
}
// hash a file or stdin
func Crc64(input *Input) (exitCode uint8) {
	return doHash(input, func() hash.Hash {
		return crc64.New(crc64.MakeTable(crc64.ISO))
	})
}
