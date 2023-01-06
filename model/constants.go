package model

const (
	OriginYesNo = "yes_no" // Origin (God) answer yes or no
	OriginOneOf = "one_of" // Origin (God) advise best variant

	ArtifactContentTypeJpeg = "image/jpeg"

	//MaxSeed = uint64(10000000000)
	MaxSeed = uint64(4294967295)

	StrategyHash  = "hash"
	StrategyScale = "scale"
)

const (
	StateError          = "error"
	StateNotWorking     = "not_working"
	StateMakingSpell    = "making_spell"
	StateMakingArtifact = "making_artifact"
	StateMakingRest     = "enjoying the result"
)
