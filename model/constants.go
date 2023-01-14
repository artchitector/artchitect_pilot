package model

const (
	// Version0, Version01, Version1, Version11
	// Artchitect have several available now versions of card generation

	Version0  = "v0"   // no tags, no StableDiffusion (first test images "Allah")
	Version01 = "v0.1" // initial tags (when I just started)
	Version1  = "v1"   // initial set of tags + InvokeAI + StableDiffusion v1.5
	Version11 = "v1.1" // more tags + InvokeAI + StableDiffusion v1.5
)

var AvailableVersions = []string{Version1, Version11}

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
	StateMakingRest     = "enjoying_the_result"
)
