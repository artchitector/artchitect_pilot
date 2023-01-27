package model

const (
	// Version0, Version01, Version1, Version11
	// Artchitect have several available now versions of card generation

	Version0  = "v0"   // no tags, no StableDiffusion (first test images "Allah")
	Version01 = "v0.1" // initial tags (when I just started)
	Version1  = "v1"   // initial set of tags + InvokeAI + StableDiffusion v1.5
	Version11 = "v1.1" // more tags + InvokeAI + StableDiffusion v1.5
	Version12 = "v1.2" // once more tags + InvokeAI + StableDiffusion v1.5
)

var AvailableVersions = []string{Version1, Version11, Version12}

const (
	MaxSeed = uint(4294967295)

	SizeXF = "xf" // very large x4 resolution, will be in future
	SizeF  = "f"  // 1024x1536
	SizeM  = "m"  // 512x768
	SizeS  = "s"  // 256x384
	SizeXS = "xs" // 128x192
)
