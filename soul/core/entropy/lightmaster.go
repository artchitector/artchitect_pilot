package entropy

import (
	"context"
	"github.com/artchitector/artchitect/soul/resources"
	"github.com/rs/zerolog/log"
	"image"
)

const (
	PhaseSource = "source"
	PhaseNoise  = "noise"
)

/*
Lightmaster отслеживает состояние энтропии в текущем кадре
Еще он передаёт детализацию обработки энтропии на gate-сервер через redis. Это нужно, чтобы на клиенте был виден
постоянный процесс обработки энтропии в виде jpeg-стримов, и видно было как картинка превращается в решение.

Не каждое состояние используется в принятии решений, многие пропускаются.
*/
type Lightmaster struct {
	webcam     *resources.Webcam
	gatekeeper *Gatekeeper
	lastFrame  image.Image
}

func NewLightmaster(webcam *resources.Webcam, gatekeeper *Gatekeeper) *Lightmaster {
	return &Lightmaster{webcam, gatekeeper, nil}
}

/*
Запускается процесс считывания кадров с веб-камеры и превращение их в float64-число
*/
func (l *Lightmaster) StartEntropyReading(ctx context.Context) error {
	ch := l.webcam.GetStream(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case i := <-ch:
			if err := l.handleSingleFrame(ctx, i); err != nil {
				log.Error().Err(err).Msgf("[lightmaster] errored handle single frame")
			}
		}
	}
}

func (l *Lightmaster) handleSingleFrame(ctx context.Context, sourceImg image.Image) error {
	if err := l.gatekeeper.PushPhase(ctx, sourceImg, PhaseSource); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseSource)
		// not stop
	}
	//var noiseImage image.Image
	//noiseImage := l.sourceToNoise(sourceImg)
	//
	//if err := l.notifyGateWithImage(ctx, noiseImage, PhaseNoise); err != nil {
	//	log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoise)
	////	not stop
	//}

	return nil
}

func (l *Lightmaster) sourceToNoise(img image.Image) image.Image {
	return nil
}
