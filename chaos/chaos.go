package chaos

import (
	"math/rand/v2"

	"github.com/cohen990/exactlyOnce/logging"
)

var logger = logging.NewRoot("chaos")

type Action[T any] func() T

type failureMode interface {
	Break()
	Fix()
}

type Chaos[F Action[T], T any] struct {
	options []failureMode
}

func (chaos *Chaos[F, T]) Register(option failureMode) {
	chaos.options = append(chaos.options, option)
}

func (chaos *Chaos[F, T]) InjectChaos(doAction F) T {
	logger := logger.Child("InjectChaos")

	if rand.Float32() > 0.5 {
		logger.Info("Today, we choose violence")
		targetIndex := rand.IntN(len(chaos.options))
		target := chaos.options[targetIndex]
		target.Break()
		result := doAction()
		target.Fix()
		return result
	} else {
		logger.Info("Today, we choose peace")
		return doAction()
	}
}
