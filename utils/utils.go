package utils

import (
	"reflect"

	. "KVBridge/types"

	"github.com/mitchellh/mapstructure"
)

func OpStrToPrefHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int {
			return data, nil
		}
		raw := data.(string)
		pref, ok := opPrefLookup[OpPreferenceString(raw)]
		if !ok {
			// then leave it as it is, someonw will throw an error down the line
			return data, nil
		}
		// else return the int
		return pref, nil
	}
}

var opPrefLookup = map[OpPreferenceString]OpPreference{
	OpLocalStr:    OpLocal,
	OpMajorityStr: OpMajority,
	OpAllStr:      OpAll,
}
