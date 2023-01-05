package speller

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"github.com/pkg/errors"
)

type spellRepository interface {
	Save(ctx context.Context, spell model.Spell) (model.Spell, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint64) (uint64, error)
}

/*
Speller generates Spell (combination of painting caption, tags and seed). Painting will be created with this Spell.
*/
type Speller struct {
	spellRepository spellRepository
	origin          origin
}

func NewSpeller(spellRepository spellRepository, origin origin) *Speller {
	return &Speller{spellRepository, origin}
}

func (s *Speller) MakeSpell(ctx context.Context) (model.Spell, error) {
	spell, err := s.generateSpell(ctx)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed to generate spell")
	}
	spell, err = s.spellRepository.Save(ctx, spell)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed to save spell in repository")
	}
	return spell, nil
}

func (s *Speller) generateSpell(ctx context.Context) (model.Spell, error) {
	selection, err := s.origin.Select(ctx, model.MaxSeed)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed to get selection")
	}
	return model.Spell{
		Idea: "hello world",
		Tags: "professional majestic oil painting, volumetric lighting, dramatic lighting, (orange)0.4",
		Seed: selection,
	}, nil
}
