package cleaner

import "time"

type Platforms = map[string]PlatformConfig

type PlatformConfig struct {
	Pairs        []PairConfig
	ExpiresEvery time.Duration
	CollectEvery time.Duration
}

type PairConfig struct {
	InternalName    string
	Mapping         string
	PrecisionMargin int32
	PrecisionAmount int32
}
