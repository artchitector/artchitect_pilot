package model

const (
	OriginYesNo = "yes_no" // Origin (God) answer yes or no
	OriginOneOf = "one_of" // Origin (God) advise best variant

	StateNotWorking     = "not_working" // Artchitect is offline
	StateMakingSpell    = "making_spell"
	StateMakingArtifact = "making_artifact"
	StateMakingRest     = "making_rest"

	ArtifactContentTypeJpeg = "image/jpeg"

	//MaxSeed = uint64(10000000000)
	MaxSeed = uint64(4294967295)
)
