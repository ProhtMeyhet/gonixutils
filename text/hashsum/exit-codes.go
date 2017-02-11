package hashsum

const(
	ERROR_HASH_FUNCTION		= 20 // the hash function has returned an error and hopefully it'll be descriptive in stderr
	ERROR_INVALID_FILE_FORMAT	= 21 // the compare file format was detected as invalid (eg wrong hashfunction)
)
