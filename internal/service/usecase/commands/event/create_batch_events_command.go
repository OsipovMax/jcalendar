package event

import (
	"context"
)

type CreateBatchEventCommand struct {
	Batch []*CreateEventCommand
}

func NewCreateBatchEventCommand(_ context.Context, batch []*CreateEventCommand) *CreateBatchEventCommand {
	return &CreateBatchEventCommand{
		Batch: batch,
	}
}
