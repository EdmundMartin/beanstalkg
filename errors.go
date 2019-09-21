package beanstalkg

import "errors"

var (
	InvalidTubeName    = errors.New(`tube name exceeds 200 characters`)
	UnexpectedResponse = errors.New(`unexpected result received from command`)
	DeadlineSoon       = errors.New(`deadline soon`)
	TimedOut           = errors.New(`timed out`)
	ExpectedCRLF       = errors.New(`expected CRLF`)
	JobTooBig          = errors.New(`job too big`)
	Draining           = errors.New(`draining`)
	Buried             = errors.New(`buried`)
	NotFound           = errors.New(`not found`)
	OutOfMemory        = errors.New(`out of memory`)
	InternalError      = errors.New(`internal error`)
	BadFormat          = errors.New(`bad format`)
	UnknownCommand     = errors.New(`unknown command`)
	NotIgnored         = errors.New(`not ignored`)
)

var errorMapping = map[string]error{
	"DEADLINE_SOON\r\n":   DeadlineSoon,
	"TIMED_OUT\r\n":       TimedOut,
	"EXPECTED_CRLF\r\n":   ExpectedCRLF,
	"JOB_TOO_BIG\r\n":     JobTooBig,
	"DRAINING\r\n":        Draining,
	"BURIED\r\n":          Buried,
	"NOT_FOUND\r\n":       NotFound,
	"OUT_OF_MEMORY\r\n":   OutOfMemory,
	"INTERNAL_ERROR\r\n":  InternalError,
	"BAD_FORMAT\r\n":      BadFormat,
	"UNKNOWN_COMMAND\r\n": UnknownCommand,
	"NOT_IGNORED\r\n":     NotIgnored,
}

func stringToError(errorString string) error {
	val, ok := errorMapping[errorString]
	if ok {
		return val
	}
	return UnexpectedResponse
}
