package winrt

import "github.com/go-ole/go-ole"

// RO_INIT_TYPE enumeration (roapi.h)
// https://learn.microsoft.com/en-us/windows/win32/api/roapi/ne-roapi-ro_init_type

type RO_INIT_TYPE uint32

const (
	RO_INIT_SINGLETHREADED RO_INIT_TYPE = 0
	RO_INIT_MULTITHREADED  RO_INIT_TYPE = 1
)

func RoInitialize(thread_type RO_INIT_TYPE) error {
	return ole.RoInitialize(uint32(thread_type))
}
