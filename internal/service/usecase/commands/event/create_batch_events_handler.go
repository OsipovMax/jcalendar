package event

//type BatchCreator interface {
//	CreateEventBatch(ctx context.Context, e *eevent.Event) (uint, error)
//}
//
//type CreateEventBatchCommandHandler struct {
//	creator BatchCreator
//}
//
//func NewCreateEventBatchCommandHandler(creator Creator) CreateEventBatchCommandHandler {
//	return CreateEventBatchCommandHandler{creator: creator}
//}
//
//func (ch *CreateEventBatchCommandHandler) Handle(ctx context.Context, command *CreateBatchEventCommand) error {
//	evs := make([]*eevent.Event, len(command.Batch), len(command.Batch))
//	for idx, subcmd := range command.Batch {
//		evs[idx] = eevent.NewEvent(
//			ctx,
//			subcmd.From,
//			subcmd.Till,
//			subcmd.CreatorID,
//			subcmd.ParticipantsIDs,
//			subcmd.Details,
//			subcmd.IsPrivate,
//			subcmd.IsRepeat,
//		)
//	}
//
//	_, err := ch.creator.CreateEvent(ctx, e)
//	if err != nil {
//		return 0, err
//	}
//
//	return e.ID, nil
//}
