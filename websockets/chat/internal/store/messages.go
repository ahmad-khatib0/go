package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

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
