package state

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/artchitector/artchitect.git/gate/model"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"image/jpeg"
	"time"
)

type Retriever struct {
	logger             zerolog.Logger
	paintingRepository paintingRepository
	decisionRepository decisionRepository
	stateRepository    stateRepository
}

func NewRetriever(
	logger zerolog.Logger,
	paintingRepository paintingRepository,
	decisionRepository decisionRepository,
	stateRepository stateRepository,
) *Retriever {
	return &Retriever{logger, paintingRepository, decisionRepository, stateRepository}
}

func (r *Retriever) CollectState(ctx context.Context) (model.CurrentState, error) {
	lastPainting, found, err := r.paintingRepository.GetLastPainting(ctx)
	if err != nil {
		return model.CurrentState{}, errors.Wrap(err, "failed to get last painting from repo")
	}
	var lastPaintingState model.LastPainting
	if found {
		lastPaintingState = model.LastPainting{ID: lastPainting.ID, Caption: lastPainting.Caption}
	} else {
		lastPaintingState = model.LastPainting{ID: 0}
	}

	lastDecision, err := r.getLastDecision(ctx)
	if err != nil {
		return model.CurrentState{}, errors.Wrap(err, "failed to get last decision")
	}

	return model.CurrentState{
		CurrentState: r.GetLastState(ctx),
		LastPainting: &lastPaintingState,
		LastDecision: lastDecision,
	}, nil
}

func (r *Retriever) GetPaintingData(ctx context.Context, paintingID uint) ([]byte, error) {
	painting, found, err := r.paintingRepository.GetPainting(ctx, paintingID)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to get painting id=%d", paintingID)
	} else if !found {
		return []byte{}, errors.Errorf("not found painting id=%d", painting)
	} else {
		return painting.Bytes, nil
	}
}

func (r *Retriever) getLastDecision(ctx context.Context) (*model.LastDecision, error) {
	ld, err := r.decisionRepository.GetLastDecision(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		return nil, nil
	}
	buf := bytes.NewBuffer(ld.Artifact)
	// now only jpeg supported
	img, err := jpeg.Decode(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode jpeg")
	}

	img = resize.Resize(0, 20, img, resize.NearestNeighbor)

	buf = new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, errors.Wrap(err, "failed to encode large jpeg")
	}

	return &model.LastDecision{
		Result: ld.Output,
		Cdate:  ld.CreatedAt,
		Image:  base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, err
}

func (r *Retriever) GetLastState(ctx context.Context) model.CurrentStateStr {
	state, err := r.stateRepository.GetLastState(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return model.CurrentStateError
	}

	now := time.Now()
	if state.CreatedAt.Before(now.Add(-1 * time.Minute)) {
		log.Warn().Msgf("too old state. %+v (id=%d)", state.CreatedAt, state.ID)
		return model.CurrentStateNotWorking
	} else {
		return model.CurrentStateStr(state.State)
	}
}
