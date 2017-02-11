package hashsum

// parse me!
func ParseType(input string) Type {
	switch input {
	case "sha512", "sha512sum":
		return SHA512
	case "md", "md5", "md5sum":
		return MD5
	case "sha1", "sha1sum":
		return SHA1
	case "sha", "sha256", "sha256sum":
		return SHA256
	case "sha224", "sha224sum":
		return SHA224
	case "sha384", "sha384sum":
		return SHA384
	case "sha512a", "sha512asum", "sha512_224", "sha512224":
		return SHA512_224
	case "sha512b", "sha512bsum", "sha512_256", "sha512256":
		return SHA512_256
	case "adler", "adler32", "adler32sum":
		return ADLER32
	case "crc", "crc32", "crc32sum":
		return CRC32
	case "crc64", "crc64sum":
		return CRC64
	}

	return NONE
}
