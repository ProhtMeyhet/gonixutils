package hashsum

const(
	NONE	Type	= 0
	SHA1	Type	= 1
	SHA224	Type	= 2
	SHA256	Type	= 3
	SHA384	Type	= 4
	SHA512_224 Type	= 5
	SHA512_256 Type	= 6
	SHA512	Type	= 7
	MD5	Type	= 8
	ADLER32 Type	= 9
	CRC32	Type	= 10
	CRC64	Type	= 11

//	FNV1	Type	= 12
//	FNV1a	Type	= 13

// XXX when adding or removing a hashtype, remember to update Type.Valid() XXX
)
