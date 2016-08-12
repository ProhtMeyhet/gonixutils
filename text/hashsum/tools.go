package hash

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

// string me!
func TypeToString(hash Type) string {
	switch(hash) {
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
	}

	return "INVALID"
}

// return the length of the hex encoded hash in byte
func HashLen(hash Type) uint {
	switch(hash) {
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
