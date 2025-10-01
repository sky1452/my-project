package models

import "time"

type module struct {
	id              rune
	moduleNum       int8         //1-3
	moduleDate      [2]time.Time // start date - end date
	attestationDate [2]time.Time // start attestation date - end attestation date
}
