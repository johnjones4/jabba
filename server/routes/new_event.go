package routes

import (
	"context"
	statusEngine "main/status"
	"main/store"
	"time"

	"github.com/johnjones4/Jabba/core"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func NewEventUseCase(s store.Store, se statusEngine.StatusEngine) usecase.IOInteractor {
	return usecase.NewIOI(new(core.Event), new(core.Event), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*core.Event)
			out = output.(*core.Event)
		)

		if in.Created.Unix() == 0 {
			in.Created = time.Now().UTC()
		}

		err := s.SaveEvent(in)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		_, err = se.HandleNewEvent(*in)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*out = *in

		return nil
	})
}
