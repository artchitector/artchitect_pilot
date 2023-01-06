package model

import "gorm.io/gorm"

// Spell - is text command to make an artwork. Spell is a combination of picture caption, tags and seed.
// Finally, Spell used by artist to make a picture.
type Spell struct {
	gorm.Model
	Tags string // additional tags to paint the picture (https://www.reddit.com/r/StableDiffusion/comments/y649yn/prompts_modifiers_to_get_midjourney_style_in/)
	Seed uint64 // specified seed (seed is from 0 to 10 000 000 000)
}
