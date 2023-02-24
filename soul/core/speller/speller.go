package speller

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const MaxTags = 28

type notifier interface {
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

type spellRepository interface {
	Save(ctx context.Context, spell model.Spell) (model.Spell, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

/*
Speller generates Spell (combination of painting caption, tags and seed). Card will be created with this Spell.
*/
type Speller struct {
	spellRepository spellRepository
	origin          origin
	notifier        notifier
	dictionaries    map[string][]string
}

func NewSpeller(spellRepository spellRepository, origin origin, notifier notifier) *Speller {
	return &Speller{spellRepository, origin, notifier, make(map[string][]string)}
}

func (s *Speller) MakeSpell(ctx context.Context, artistState *model.CreationState) (model.Spell, error) {
	spell, err := s.generateSpell(ctx, artistState)
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

func (s *Speller) generateSpell(ctx context.Context, state *model.CreationState) (model.Spell, error) {
	version, err := s.selectVersion(ctx)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed select version")
	}
	state.Version = version
	s.notify(ctx, state)
	selection, err := s.origin.Select(ctx, model.MaxSeed)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed to get selection")
	}
	state.Seed = selection
	s.notify(ctx, state)
	tags, err := s.generateTags(ctx, version, state)
	if err != nil {
		return model.Spell{}, errors.Wrap(err, "[speller] failed generate tags")
	}
	return model.Spell{
		Tags:    strings.Join(tags, ","),
		Seed:    selection,
		Version: version,
	}, nil
}

func (s *Speller) generateTags(ctx context.Context, version string, state *model.CreationState) ([]string, error) {
	dictionary, err := s.getDictionary(ctx, version)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get Dictionary")
	}

	tagsToTake, err := s.origin.Select(ctx, MaxTags)
	tagsToTake += 1 // Select returns [0,MaxTags). Plus 1 is ok

	tags := make([]string, 0, tagsToTake)
	if err != nil {
		return []string{}, errors.Wrap(err, "[speller][generateTags] failed get tagsToTake")
	}
	state.TagsCount = tagsToTake
	s.notify(ctx, state)

	allowedTagsLen := uint(len(dictionary))
	for i := uint(0); i < tagsToTake; i++ {
		tag, err := s.origin.Select(ctx, allowedTagsLen)
		if err != nil {
			return []string{}, errors.Wrap(err, "[speller][generateTags] failed get tag number")
		}
		tags = append(tags, dictionary[tag])
		state.Tags = append(state.Tags, dictionary[tag])
		s.notify(ctx, state)
	}
	return tags, nil
}

func (s *Speller) getDictionary(ctx context.Context, version string) ([]string, error) {
	dict, found := s.dictionaries[version]
	if found {
		return dict, nil
	}

	var filename string
	switch version {
	case model.Version1:
		filename = "files/tags_v1.yaml"
	case model.Version11:
		filename = "files/tags_v11.yaml"
	case model.Version12:
		filename = "files/tags_v12.yaml"
	case model.Version20:
		filename = "files/tags_v2.yaml"

	default:
		return []string{}, errors.Errorf("[speller] unknown version %s. failed to load file", version)
	}

	log.Info().Msgf("[speller] loading tags file %s", filename)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to load yaml file")
	}
	tags := []string{}
	err = yaml.Unmarshal(yamlFile, &tags)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to parse yaml file")
	}
	log.Info().Msgf("[speller] version=%s, file=%s, loaded tags: %d. First ten: %s", version, filename, len(tags), strings.Join(tags[0:10], ", "))
	s.dictionaries[version] = tags

	return s.dictionaries[version], nil
}

func (s *Speller) selectVersion(ctx context.Context) (string, error) {
	count := len(model.AvailableVersions)
	idx, err := s.origin.Select(ctx, uint(count))
	return model.AvailableVersions[idx], errors.Wrap(err, "[speller] failed to select version from origin")
}

func (s *Speller) notify(ctx context.Context, state *model.CreationState) {
	if err := s.notifier.NotifyCreationState(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[speller] failed notify artist state")
	}
}
