package abstract

const(
	/*
	*  SUCCESS is anything < 10, except, for historic reasons, 1
	*  but >= 1 is only partly successful
	*  but >= 9 empty or partly empty
	*/

	// Aus dem Hintergrund müsste Rahn schießen, Rahn schießt - TOR, TOR, TOR, TOR!
	SUCCESS		= 0

	// for historic reasons, say 1 is FAILED
	FAILED		= 1

	// from a given list, at least one part failed. check stderr.
	PARTLY		= 2


	/*
	*  SUCCESS values 3 - 8 are reserved for individual programs
	*/


	// output is empty, so a new line was not printed.
	EMPTY_OUTPUT	= 9

	// some of the output wasn't outputted because it was empty.
	PARTLY_EMPTY	= 10

	// maximum success. used internally for if exitCode > SUCCESSFUL
	SUCCESSFUL	= 10


	/*
	*  error values 20 - 99 are reserved for individual programs
	*/


	/*
	*  anything 100-199 is a generic error with the user given parameters
	*/

	// this shouldn't happen. declare the program received an error, but doesn't
	// know how to handle it or even what it is.
	ERROR_UNSPECIFIED		= 100

	// it finally happend
	ERROR_DISK_FULL			= 101

	// access denied in resource
	ERROR_ACCESS_DENIED		= 102

	// couldn't write to resource
	ERROR_WRITE_ACCESS_DENID	= 103

	// no input given
	ERROR_NO_INPUT			= 104

	// unknown (kill) signal
	ERROR_INVALID_SIGNAL		= 105

	// error parsing input
	ERROR_PARSING			= 106

	// invalid argument
	ERROR_INVALID_ARGUMENT		= 107

	// couldn't get working directory
	ERROR_WORKING_DIRECTORY		= 108

	// doesn't exist
	ERROR_FILE_NOT_FOUND		= 109

	// file already exist. file will not be overriden
	ERROR_FILE_EXIST		= 110

	// open error
	ERROR_OPENING_FILE		= 111

	// reading error on io.Read
	ERROR_READING			= 112

	// couldn't stat inode
	ERROR_STAT			= 113

	// TODO remove
	ERROR_COMPAT			= 199



	/*
	*  anything >= 200 is an internal error that might or might not have anything to do with
	*  user input
	*/

	// an unknown or unhandled internal error occured.
	ERROR_INTERNAL			= 200

	// access denied to internally used resource. config file for example.
	ERROR_INTERNAL_ACCESS_DENIED	= 201

	// write access denied to internally used resource. config file for example.
	ERROR_INTERNAL_WRITE_ACCESS_DENIED = 202

	// unhandled error. naughty programmer!
	ERROR_UNHANDLED			= 254

	// this error should'nt be used. please be more fine grained.
	ERROR_GENERIC			= 255
)
