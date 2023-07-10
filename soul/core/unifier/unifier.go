package unifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

const Thumb100Size = 4
const Thumb1000Size = 6
const Thumb10000Size = 8
const UpdatePeriod10000 = 100 // Единство 10к обновляется каждые 100 карточек
const UpdatePeriod1000 = 50   // Единство 1к обновляется каждые 50 карточек
const UpdatePeriod100 = 10    // Единство 100 обновляется каждые 10 карточек

type unityRepository interface {
	GetUnity(mask string) (model.Unity, error)
	CreateUnity(mask string) (model.Unity, error)
	CreateUnityByCard(cardID uint, rank uint) (model.Unity, error)
	SaveUnity(unity model.Unity) (model.Unity, error)
	GetNextUnityForWork() (model.Unity, error)
	GetUnityByCard(cardID uint, rank uint) (model.Unity, error)
}

type cardRepository interface {
	GetAnyCardIDFromHundred(ctx context.Context, rank uint, start uint) (uint, error)
	GetPreviousCardID(ctx context.Context, cardID uint) (uint, error)
	GetCard(ctx context.Context, ID uint) (model.Card, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

type combinator interface {
	CombineThumb(ctx context.Context, cardIDs []uint, mask string, version int) error
}

type notifier interface {
	NotifyUnity(ctx context.Context, state model.UnityState) error
}

type artchitectBot interface {
	SendUnityTo10Min(ctx context.Context, unity model.Unity) error
}

type Unifier struct {
	unityRepository unityRepository
	cardRepository  cardRepository
	origin          origin
	combinator      combinator
	notifier        notifier
	artchitectBot   artchitectBot
}

func NewUnifier(unityRepository unityRepository, cardRepository cardRepository, origin origin, combinator combinator, notifier notifier, artchitectBot artchitectBot) *Unifier {
	return &Unifier{unityRepository, cardRepository, origin, combinator, notifier, artchitectBot}
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

func (u *Unifier) UpdateUnitiesByNewCard(ctx context.Context, cardID uint) (bool, error) {
	// Когда создаётся новая картина, то в это время нужно обновить единства, в которых она состоит.
	// Каждая картина входит в единство сотня, тысяча, десятитысяча (а в будущем и для 100к тоже будет единство).
	// Процесс обновления единств тяжелый, поэтому обновление не каждую карточку
	// Для десятитысячного единства каждые 100 карточек, для тысячного единства 50 карточек, а для сотни каждые 10 карточек
	// Для понимания перехода через 100, 50 и 10 карточек берётся ID предыдущей карточки.
	// Это для неё перестраивается единство, а не для текущей
	prevCardID, err := u.cardRepository.GetPreviousCardID(ctx, cardID)
	if err != nil {
		return false, errors.Wrapf(err, "[unifier] failed to GetPreviousCardID for %d", cardID)
	}

	worked := false

	// working with update of 10k unity
	if int(math.Floor(float64(prevCardID/UpdatePeriod10000))) != int(math.Floor(float64(cardID/UpdatePeriod10000))) {
		log.Info().Msgf("[unifier] there is %d-class change between cards prev:%d and new:%d", model.Rank10000, prevCardID, cardID)
		// need update 10k unity
		if un, err := u.unityRepository.GetUnityByCard(prevCardID, model.Rank10000); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.Wrapf(err, "[unifier] failed to find unity-%d for card %d", model.Rank10000, cardID)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// just create new unity
			un, err = u.unityRepository.CreateUnityByCard(prevCardID, model.Rank10000)
			if err != nil {
				return false, errors.Wrapf(err, "[unifier] failed to find unity-%d for card %d", model.Rank10000, cardID)
			}
			log.Info().Msgf("[unifier] created unity %s/%d for card %d", un.Mask, un.Rank, prevCardID)
			worked = true
		} else {
			// unity already exists
			un.State = model.UnityStateReunification
			if _, err := u.unityRepository.SaveUnity(un); err != nil {
				return false, errors.Wrapf(err, "[unifier] failed to update unity-%s for card %d", un.Mask, cardID)
			}
			log.Info().Msgf("[unifier] reunify unity %s/%d for card %d", un.Mask, un.Rank, prevCardID)
			worked = true
		}
	}

	// working with update of 1k unity
	if int(math.Floor(float64(prevCardID/UpdatePeriod1000))) != int(math.Floor(float64(cardID/UpdatePeriod1000))) {
		log.Info().Msgf("[unifier] there is %d-class change between cards prev:%d and new:%d", model.Rank1000, prevCardID, cardID)
		// need update 1k unity
		if un, err := u.unityRepository.GetUnityByCard(prevCardID, model.Rank1000); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.Wrapf(err, "[unifier] failed to find unity-%d for card %d", model.Rank1000, cardID)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// just create new unity
			un, err = u.unityRepository.CreateUnityByCard(prevCardID, model.Rank1000)
			if err != nil {
				return false, errors.Wrapf(err, "[unifier] failed to find unity-%d for card %d", model.Rank1000, cardID)
			}
			log.Info().Msgf("[unifier] created unity %s/%d for card %d", un.Mask, un.Rank, prevCardID)
			worked = true
		} else {
			// unity already exists
			un.State = model.UnityStateReunification
			if _, err := u.unityRepository.SaveUnity(un); err != nil {
				return false, errors.Wrapf(err, "[unifier] failed to update unity-%s for card %d", un.Mask, cardID)
			}
			log.Info().Msgf("[unifier] reunify unity %s/%d for card %d", un.Mask, un.Rank, prevCardID)
			worked = true
		}
	}

	// working with update of 100 unity
	if int(math.Floor(float64(prevCardID/UpdatePeriod100))) != int(math.Floor(float64(cardID/UpdatePeriod100))) {
		log.Info().Msgf("[unifier] there is %d-class change between cards prev:%d and new:%d", model.Rank100, prevCardID, cardID)
		// need update 1k unity
		if un, err := u.unityRepository.GetUnityByCard(prevCardID, model.Rank100); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.Wrapf(err, "[unifier] failed to find unity-%d for card %d", model.Rank100, cardID)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// just create new unity
			un, err = u.unityRepository.CreateUnityByCard(prevCardID, model.Rank100)
			if err != nil {
				return false, errors.Wrapf(err, "[unifier] failed to find unity-%d for card %d", model.Rank100, cardID)
			}
			log.Info().Msgf("[unifier] created unity %s/%d for card %d", un.Mask, un.Rank, prevCardID)
			worked = true
		} else {
			// unity already exists
			un.State = model.UnityStateReunification
			if _, err := u.unityRepository.SaveUnity(un); err != nil {
				return false, errors.Wrapf(err, "[unifier] failed to update unity-%s for card %d", un.Mask, cardID)
			}
			log.Info().Msgf("[unifier] reunify unity %s/%d for card %d", un.Mask, un.Rank, prevCardID)
			worked = true
		}
	}

	return worked, nil
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
	// Повышаем версию единства
	unity.Version = unity.Version + 1
	state.SetUnity(unity)
	u.notify(ctx, state)
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
	// завершает унификацию
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

		if child.State == model.UnityStateUnified || child.State == model.UnityStateSkipped {
			// skipped unity will only rebuild on reunification
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
	if unity.State == model.UnityStateReunification || unity.State == model.UnityStateSkipped {
		unity.Leads = ""
	}
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

			selection, err := u.origin.Select(ctx, model.Rank100)
			if err != nil {
				return model.Unity{}, errors.Wrapf(err, "[unifier] failed to get data from origin")
			}
			lead := unity.Start() + selection // Выбираем даже несуществующие карточки. Они будут заполняться чёрным цветом.
			if _, err := u.cardRepository.GetCard(ctx, lead); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				log.Info().Msgf("[unifier] not found lead card %d, use 0", lead)
				lead = 0
			} else if err != nil {
				return model.Unity{}, errors.Wrapf(err, "[unifier] failed to check card %d existence", lead)
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
	if err := u.combinator.CombineThumb(ctx, leads, unity.Mask, unity.Version); err != nil {
		log.Error().Err(err).Msgf("[unifier] failed combine thumb %s", unity.Mask)
	}
	state.SetState(model.UnityStatePrepareThumb, 1, 1)
	state.SetThumb(unity.Mask)
	u.notify(ctx, state)

	return unity, nil
}

func (u *Unifier) finishUnification(ctx context.Context, unity model.Unity, state *model.UnityState) (model.Unity, error) {
	unityState := model.UnityStateSkipped

	if unity.Leads == "" {
		log.Info().Msgf("[unifier] unity %s have no leads", unity)
		unity.State = model.UnityStateSkipped
	} else {
		var leads []uint
		if err := json.Unmarshal([]byte(unity.Leads), &leads); err != nil {
			return model.Unity{}, errors.Wrapf(err, "[unifier] failed to unmarshal leads")
		}
		for _, lead := range leads {
			if lead > 0 {
				unityState = model.UnityStateUnified
			}
		}
	}

	unity.State = unityState
	saved, err := u.unityRepository.SaveUnity(unity)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[unifier] failed finish unity %s", unity.Mask)
	}

	isCompleted, err := u.isUnityCompleted(unity)
	if (u.artchitectBot != nil && unity.Rank == model.Rank10000 && (unity.Version%5 == 0 || isCompleted)) ||
		(unity.Rank == model.Rank1000 && isCompleted) {
		if err := u.artchitectBot.SendUnityTo10Min(ctx, unity); err != nil {
			log.Error().Err(err).Msgf("[unifier] failed to notify bot about new unity %s-%d", unity.Mask, unity.Version)
		}
	}

	for i := 1; i <= 10; i++ {
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
	log.Info().Msgf("[unifier] finished unity %s with state %s with version %d", unity, unity.State, unity.Version)
	return saved, nil
}

func (u *Unifier) notify(ctx context.Context, state *model.UnityState) {
	if err := u.notifier.NotifyUnity(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[unifier] failed notify state")
	}
}

func (u *Unifier) isUnityCompleted(unity model.Unity) (bool, error) {
	var leads []uint
	if err := json.Unmarshal([]byte(unity.Leads), &leads); err != nil {
		return false, errors.Wrapf(err, "[unifier] failed to unmarshal leads %s", unity.Mask)
	}
	for _, lead := range leads {
		if lead == 0 {
			// not completed, if there is 0 lead
			return false, nil
		}
	}
	return true, nil
}
