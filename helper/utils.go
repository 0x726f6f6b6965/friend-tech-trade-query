package helper

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// IsValidAddress - validate hex address
func IsValidAddress(addr interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := addr.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsValidTx - validate hex transaction
func IsValidTx(tx interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	switch v := tx.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// Empty - check string is empty
func Empty(s string) bool {
	return strings.Trim(s, " ") == ""
}

func GeneralDuration(times, minTime, maxTime int, duration time.Duration) time.Duration {
	return time.Duration(times+RandInt(minTime, maxTime)) * duration
}
func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
