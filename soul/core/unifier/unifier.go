package unifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

const Thumb100Size = 4
const Thumb1000Size = 6
const Thumb10000Size = 8

type unityRepository interface {
	GetUnity(mask string) (model.Unity, error)
	CreateUnity(mask string) (model.Unity, error)
	SaveUnity(unity model.Unity) (model.Unity, error)
	GetNextUnityForWork() (model.Unity, error)
}

type cardRepository interface {
	GetAnyCardIDFromHundred(ctx context.Context, rank uint, start uint) (uint, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

type combinator interface {
	CombineThumb(ctx context.Context, cardIDs []uint, mask string) error
}

type notifier interface {
	NotifyUnity(ctx context.Context, state model.UnityState) error
}

type Unifier struct {
	unityRepository unityRepository
	cardRepository  cardRepository
	origin          origin
	combinator      combinator
	notifier        notifier
}

func NewUnifier(unityRepository unityRepository, cardRepository cardRepository, origin origin, combinator combinator, notifier notifier) *Unifier {
	return &Unifier{unityRepository, cardRepository, origin, combinator, notifier}
}

func (u *Unifier) WorkOnce(ctx context.Context) (bool, error) {
	un, err := u.unityRepository.GetNextUnityForWork()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.Wrapf(err, "[unifier] failed get unity for work")
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	state := model.UnityState{}
	_, err = u.Unify(ctx, un, &state)
	if err != nil {
		return false, errors.Wrapf(err, "[unifier] failed to Unify %s", un.Mask)
	}
	return true, nil
}

/*
см. model/unity.go
Унификатор создаёт единства.
На вход получает подготовленное единство с заполненной маской и создаёт это единство (заполняет его карточками)
У каждого единства своя карточка, состоящая из набора вложенных карточек.
*/
func (u *Unifier) Unify(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	log.Info().Msgf("[unifier] unify unity %s", unity)

	state.Add(unity)
	u.notify(ctx, state)

	var err error
	// Проверим, что все дочерние unity уже созданы в БД
	if unity, err = u.checkChildrenExists(ctx, unity, state); err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed check children exists for %s", unity)
	}
	select {
	case <-ctx.Done():
		return model.Unity{}, nil
	default:
	}
	// Постепенно объединяем все дочерние единства, если они еще не объединены
	if unity, err = u.fillChildren(ctx, unity, state); err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed to fill children for %s", unity)
	}
	select {
	case <-ctx.Done():
		return model.Unity{}, nil
	default:
	}
	// Формируем список лидеров данного единства
	if unity, err = u.promoteLeads(ctx, unity, state); err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed to promote leads for %s", unity)
	}
	select {
	case <-ctx.Done():
		return model.Unity{}, nil
	default:
	}
	// Формирует картинку
	if unity, err = u.saveThumb(ctx, unity, state); err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed to save thumb for %s", unity)
	}
	select {
	case <-ctx.Done():
		return model.Unity{}, nil
	default:
	}
	if unity, err = u.finishUnification(ctx, unity, state); err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed to finish unity %s", unity)
	}

	state.Remove()
	u.notify(ctx, state)

	return unity, nil
}

func (u *Unifier) checkChildrenExists(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	if unity.Rank == model.Rank100 { // chilren of 100 is just cards
		return unity, nil
	}
	var children []model.Unity
	for i := 0; i < 10; i++ {
		submask := strings.Replace(unity.Mask, "X", fmt.Sprintf("%d", i), 1)
		child, err := u.unityRepository.GetUnity(submask)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Unity{}, errors.Wrapf(err, "[unifier] get child unity %s failed", submask)
		} else if err == nil {
			// child already exists
			log.Info().Msgf("[unifier] child %s already exists in %s", submask, unity.Mask)
			children = append(children, child)
			state.AddChild(child)
			state.SetState(model.UnityStateCollectingChildren, i+1, 10)
			u.notify(ctx, state)
		} else {
			child, err = u.unityRepository.CreateUnity(submask)
			if err != nil {
				return model.Unity{}, errors.Wrapf(err, "[unifier] create child unity %s failed", submask)
			}
			log.Info().Msgf("[unifier] child %s created in %s", submask, unity.Mask)
			children = append(children, child)
			state.AddChild(child)
			state.SetState(model.UnityStateCollectingChildren, i+1, 10)
			u.notify(ctx, state)
		}
	}
	unity.Children = children
	return unity, nil
}

func (u *Unifier) fillChildren(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	if unity.Rank == model.Rank100 { // children of 100 is just cards
		return unity, nil
	}
	for idx, child := range unity.Children {
		select {
		case <-ctx.Done():
			return model.Unity{}, nil
		default:
		}

		state.SetState(model.UnityStateUnifyChildren, idx+1, len(unity.Children))
		u.notify(ctx, state)

		if child.State == model.UnityStateUnified {
			continue
		}
		log.Info().Msgf("[unifier] unify child %s", child.Mask)
		child, err := u.Unify(ctx, child, state)
		if err != nil {
			return model.Unity{}, errors.Wrapf(err, "[unifier] unified child %s", child.Mask)
		}
		unity.Children[idx] = child

		state.SetChildState(child.Mask, child.State)
		u.notify(ctx, state)
	}
	return unity, nil
}

func (u *Unifier) promoteLeads(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	if unity.Rank == model.Rank100 {
		leadsCount := Thumb100Size * Thumb100Size
		if unity.Leads != "" {
			// load old saved leads
			var oldLeads []uint
			if err := json.Unmarshal([]byte(unity.Leads), &oldLeads); err != nil {
				return model.Unity{}, errors.Wrapf(err, "[unifier] failed to unmarshal leads of %s", unity.Mask)
			}
			if len(oldLeads) == leadsCount {
				// leads already selected

				state.SetLeads(oldLeads)
				state.SetState(model.UnityStatePromoteLeads, len(oldLeads), leadsCount)
				u.notify(ctx, state)

				log.Info().Msgf("[unifier] leads already selected for %s", unity.Mask)
				return unity, nil
			}
		}

		leads := make([]uint, 0, leadsCount)
		for i := 0; i < leadsCount; i++ {
			select {
			case <-ctx.Done():
				return model.Unity{}, nil
			default:
			}
			lead, err := u.cardRepository.GetAnyCardIDFromHundred(ctx, model.Rank100, unity.Start())
			if err != nil {
				return model.Unity{}, errors.Wrapf(err, "[unifier] failed to get card for lead. i=%d", i)
			}
			log.Info().Msgf("[unifier] selected lead %d for %s", lead, unity.Mask)
			leads = append(leads, lead)
			state.AddLead(lead)
			state.SetState(model.UnityStatePromoteLeads, len(leads), leadsCount)
			u.notify(ctx, state)
		}
		leadsj, err := json.Marshal(leads)
		if err != nil {
			return model.Unity{}, errors.Wrapf(err, "[unifier] failed to marshal leads")
		}
		unity.Leads = string(leadsj)
	} else {
		var currentLeaders []uint
		var allLeaders []uint
		for _, child := range unity.Children {
			var subleaders []uint
			if err := json.Unmarshal([]byte(child.Leads), &subleaders); err != nil {
				return model.Unity{}, errors.Wrapf(err, "[unifier] failed to unmarshal subleads %s", child.Leads)
			}
			allLeaders = append(allLeaders, subleaders...)
		}
		leadersCount := 0
		if unity.Rank == model.Rank1000 {
			leadersCount = Thumb1000Size * Thumb1000Size
		} else if unity.Rank == model.Rank10000 {
			leadersCount = Thumb10000Size * Thumb10000Size
		}
		if leadersCount == 0 {
			return model.Unity{}, errors.Errorf("[unifier] wrong rank %d", unity.Rank)
		}
		for len(currentLeaders) < leadersCount {
			selection, err := u.origin.Select(ctx, uint(len(allLeaders)))
			if err != nil {
				return model.Unity{}, errors.Errorf("[unifier] failed get selection from origin. all lead count: %d", len(allLeaders))
			}
			selectedLeader := allLeaders[selection]
			currentLeaders = append(currentLeaders, selectedLeader)
			state.AddLead(selectedLeader)
			state.SetState(model.UnityStatePromoteLeads, len(currentLeaders), leadersCount)
			u.notify(ctx, state)
		}
		leadsj, err := json.Marshal(currentLeaders)
		if err != nil {
			return model.Unity{}, errors.Wrapf(err, "[unifier] failed to marshal currentLeaders")
		}
		unity.Leads = string(leadsj)
	}

	unity, err := u.unityRepository.SaveUnity(unity)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed save unity %s", unity.Mask)
	}

	state.SetUnity(unity)
	u.notify(ctx, state)

	log.Info().Msgf("[unifier] saved leads for unity %s", unity.Mask)
	return unity, nil
}

func (u *Unifier) saveThumb(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	log.Info().Msgf("[unifier] save thumb of %s", unity.Mask)

	select {
	case <-ctx.Done():
		return model.Unity{}, nil
	default:
	}

	state.SetState(model.UnityStatePrepareThumb, 0, 1)
	u.notify(ctx, state)

	var leads []uint
	if err := json.Unmarshal([]byte(unity.Leads), &leads); err != nil {
		log.Error().Err(err).Msgf("[unifier] failed unmarshal unity %s leads %s", unity.Mask, unity.Leads)
	}
	if err := u.combinator.CombineThumb(ctx, leads, unity.Mask); err != nil {
		log.Error().Err(err).Msgf("[unifier] failed combine thumb %s", unity.Mask)
	}
	state.SetState(model.UnityStatePrepareThumb, 1, 1)
	state.SetThumb(unity.Mask)
	u.notify(ctx, state)

	return unity, nil
}

func (u *Unifier) finishUnification(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	unity.State = model.UnityStateUnified
	saved, err := u.unityRepository.SaveUnity(unity)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed finish unity %s", unity.Mask)
	}
	for i := 1; i <= 2; i++ {
		select {
		case <-ctx.Done():
			return model.Unity{}, nil
		default:
		}
		state.SetUnity(saved)
		state.SetState(model.UnityStateFinish, i, 10)
		u.notify(ctx, state)
		time.Sleep(time.Second)
	}
	log.Info().Msgf("[unifier] finished unity %s", unity.Mask)
	return saved, nil
}

func (u *Unifier) notify(ctx context.Context, state *model.UnityState) {
	if err := u.notifier.NotifyUnity(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[unifier] failed notify state")
	}
}
