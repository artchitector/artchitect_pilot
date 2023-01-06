package speller

import (
	"context"
	model2 "github.com/artchitector/artchitect.git/model"
	"github.com/artchitector/artchitect.git/soul/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const MaxTags = 28

type spellRepository interface {
	Save(ctx context.Context, spell model.Spell) (model.Spell, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint64, saveDecision bool) (uint64, error)
}

/*
Speller generates Spell (combination of painting caption, tags and seed). Painting will be created with this Spell.
*/
type Speller struct {
	spellRepository spellRepository
	origin          origin
	dictionary      []string
}

func NewSpeller(spellRepository spellRepository, origin origin) *Speller {
	return &Speller{spellRepository, origin, []string{}}
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
	log.Info().Msgf("speller made spell: %+v", spell)
	return spell, nil
}

func (s *Speller) generateSpell(ctx context.Context) (model.Spell, error) {
	selection, err := s.origin.Select(ctx, model2.MaxSeed, true)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed to get selection")
	}
	tags, err := s.generateTags(ctx)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed generate tags")
	}
	return model.Spell{
		Tags: strings.Join(tags, ","),
		Seed: selection,
	}, nil
}

func (s *Speller) generateTags(ctx context.Context) ([]string, error) {
	dictionary, err := s.getDictionary(ctx)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get Dictionary")
	}

	tagsToTake, err := s.origin.Select(ctx, MaxTags, false)
	tags := make([]string, 0, tagsToTake)
	if err != nil {
		return []string{}, errors.Wrap(err, "[speller][generateTags] failed get tagsToTake")
	}

	allowedTagsLen := uint64(len(dictionary))
	for i := uint64(0); i < tagsToTake; i++ {
		tag, err := s.origin.Select(ctx, allowedTagsLen, false)
		if err != nil {
			return []string{}, errors.Wrap(err, "[speller][generateTags] failed get tag number")
		}
		tags = append(tags, dictionary[tag])
	}
	return tags, nil
}

func (s *Speller) getDictionary(ctx context.Context) ([]string, error) {
	if len(s.dictionary) == 0 {
		yamlFile, err := os.ReadFile("files/tags.yaml")
		if err != nil {
			return []string{}, errors.Wrap(err, "failed to load yaml file")
		}
		tags := []string{}
		err = yaml.Unmarshal(yamlFile, &tags)
		if err != nil {
			return []string{}, errors.Wrap(err, "failed to parse yaml file")
		}
		log.Info().Msgf("[speller] loaded tags: %d. First ten: %s", len(tags), strings.Join(tags[0:10], ", "))
		s.dictionary = tags
	}

	return s.dictionary, nil
}
