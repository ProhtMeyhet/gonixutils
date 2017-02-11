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
)

// hash a file or stdin
func Hashsum(input *Input) (exitCode uint8) {
	if input.Compare {
		return Compare(input)
	}

	return hashsum(input, input.Type.Factory(input))
}

// hash a file or stdin
func Md5(input *Input) (exitCode uint8) {
	return hashsum(input, md5.New)
}

// hash a file or stdin
func Sha1(input *Input) (exitCode uint8) {
	return hashsum(input, sha1.New)
}

// hash a file or stdin
func Sha224(input *Input) (exitCode uint8) {
	return hashsum(input, sha256.New224)
}

// hash a file or stdin
func Sha256(input *Input) (exitCode uint8) {
	return hashsum(input, sha256.New)
}

// hash a file or stdin
func Sha384(input *Input) (exitCode uint8) {
	return hashsum(input, sha512.New384)
}

// hash a file or stdin
func Sha512(input *Input) (exitCode uint8) {
	return hashsum(input, sha512.New)
}

// hash a file or stdin
func Sha512_224(input *Input) (exitCode uint8) {
	return hashsum(input, sha512.New512_224)
}
// hash a file or stdin
func Sha512_256(input *Input) (exitCode uint8) {
	return hashsum(input, sha512.New512_256)
}
// hash a file or stdin
func Adler32(input *Input) (exitCode uint8) {
	return hashsum(input, func() hash.Hash {
		return adler32.New()
	})
}
// hash a file or stdin
func Crc32(input *Input) (exitCode uint8) {
	return hashsum(input, func() hash.Hash {
		return crc32.NewIEEE()
	})
}
// hash a file or stdin
func Crc64(input *Input) (exitCode uint8) {
	return hashsum(input, func() hash.Hash {
		return crc64.New(crc64.MakeTable(crc64.ISO))
	})
}
