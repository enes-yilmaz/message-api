package utils

import (
	"MESSAGEAPI/src/pkg/errors"
	"strconv"
)

func ValidateLimitOffset(limit, offset string, defaultLimit int) (int, int, error) {
	var lmt int
	var offst int

	if len(limit) > 0 {

		if l, err := strconv.Atoi(limit); l <= 0 && err == nil {
			return 0, 0, errors.LimitLessThanZeroError
		}

		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return 0, 0, errors.LimitParsingError
		}

		if l > 300 {
			l = 300
		}

		lmt = int(l)

	} else {
		lmt = defaultLimit
	}

	if len(offset) > 0 {
		o, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return 0, 0, errors.OffsetParsingError
		}
		if o < 0 {
			return 0, 0, errors.OffsetLessThanZeroError
		}
		offst = int(o)
	}

	return lmt, offst, nil
}
