package store

import (
	"sort"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

func (s *Store) MsgSave(msg *types.Message, attachmentURLs []string, readBySender bool) (error, bool) {
	msg.InitTimes()
	msg.SetUid(s.UidGen.Get())

	// Increment topic's or user's SeqId
	err := s.adp.Topics().UpdateOnMessage(msg.Topic, msg)
	if err != nil {
		return err, false
	}

	err = s.adp.Messages().Save(msg)
	if err != nil {
		return err, false
	}

	markedReadBySender := false
	// Mark message as read by the sender.
	if readBySender {
		// Make sure From is valid, otherwise we will reset values for all subscribers.
		fromUid := types.ParseUid(msg.From)
		if !fromUid.IsZero() {
			// Ignore the error here. It's not a big deal if it fails.
			if subErr := s.adp.Subscriptions().Update(
				msg.Topic,
				fromUid,
				map[string]interface{}{"RecvSeqId": msg.SeqId, "ReadSeqId": msg.SeqId},
			); subErr != nil {
				s.logger.Sugar().Warnf("topic[%s]: failed to mark message (seq: %d) read by sender - err: %+v", msg.Topic, msg.SeqId, subErr)
			} else {
				markedReadBySender = true
			}
		}
	}

	if len(attachmentURLs) > 0 {
		var attachments []string
		for _, url := range attachmentURLs {
			// Convert attachment URLs to file IDs.
			if fid := s.mediaHandler.GetIdFromUrl(url); !fid.IsZero() {
				attachments = append(attachments, fid.String())
			}
		}
		if len(attachments) > 0 {
			return s.adp.Files().LinkAttachments("", types.ZeroUid, msg.Uid(), attachments), markedReadBySender
		}
	}

	return nil, markedReadBySender
}

// DeleteList deletes multiple messages defined by a list of ranges.
func (s *Store) MsgDeleteList(topic string, delID int, forUser types.Uid, ranges []types.Range) error {
	var toDel *types.DelMessage
	if delID > 0 {
		toDel = &types.DelMessage{
			Topic:       topic,
			DelId:       delID,
			DeletedFor:  forUser.String(),
			SeqIdRanges: ranges}
		toDel.SetUid(s.UidGen.Get())
		toDel.InitTimes()
	}

	err := s.adp.Messages().DeleteList(topic, toDel)
	if err != nil {
		return err
	}

	// TODO: move to adapter.
	if delID > 0 {
		// Record ID of the delete transaction
		err = s.adp.Topics().Update(topic, map[string]interface{}{"DelId": delID})
		if err != nil {
			return err
		}

		// Soft-deleting will update one subscription, hard-deleting will ipdate all.
		// Soft- or hard- is defined by the forUser being defined.
		err = s.adp.Subscriptions().Update(topic, forUser, map[string]interface{}{"DelId": delID})
		if err != nil {
			return err
		}
	}

	return err
}

// GetAll returns multiple messages.
func (s *Store) MsgGetAll(topic string, forUser types.Uid, opt *types.QueryOpt) ([]types.Message, error) {
	return s.adp.Messages().GetAll(topic, forUser, opt)
}

// GetDeleted returns the ranges of deleted messages and the largest DelId reported in the list.
func (s *Store) MsgGetDeleted(topic string, forUser types.Uid, opt *types.QueryOpt) ([]types.Range, int, error) {
	dmsgs, err := s.adp.Messages().GetDeleted(topic, forUser, opt)
	if err != nil {
		return nil, 0, err
	}

	var ranges []types.Range
	var maxID int
	// Flatten out the ranges
	for i := range dmsgs {
		dm := &dmsgs[i]
		if dm.DelId > maxID {
			maxID = dm.DelId
		}
		ranges = append(ranges, dm.SeqIdRanges...)
	}

	sort.Sort(types.RangeSorter(ranges))
	ranges = types.RangeSorter(ranges).Normalize()

	return ranges, maxID, nil
}
