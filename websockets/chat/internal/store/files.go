package store

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// StartUpload records that the given user initiated a file upload
func (s *Store) FilesStartUpload(fd *types.FileDef) error {
	fd.Status = types.UploadStarted
	return s.adp.Files().StartUpload(fd)
}

// FinishUpload marks started upload as successfully finished or failed.
func (s *Store) FilesFinishUpload(fd *types.FileDef, success bool, size int64) (*types.FileDef, error) {
	return s.adp.Files().FinishUpload(fd, success, size)
}

// Get fetches a file record for a unique file id.
func (s *Store) FilesGet(fid string) (*types.FileDef, error) {
	return s.adp.Files().Get(fid)
}

// DeleteUnused removes unused attachments and avatars.
func (s *Store) FilesDeleteUnused(olderThan time.Time, limit int) error {
	toDel, err := s.adp.Files().DeleteUnused(olderThan, limit)
	if err != nil {
		return err
	}

	if len(toDel) > 0 {
		s.logger.Sugar().Warnf("deleting media", toDel)
		return s.GetMediaHandler().Delete(toDel)
	}

	return nil
}

// LinkAttachments connects earlier uploaded attachments to a
//
// message or topic to prevent it from being garbage collected.
func (s *Store) FilesLinkAttachments(topic string, msgId types.Uid, attachments []string) error {
	// Convert attachment URLs to file IDs.
	var fids []string
	for _, url := range attachments {
		if fid := s.mediaHandler.GetIdFromUrl(url); !fid.IsZero() {
			fids = append(fids, fid.String())
		}
	}

	if len(fids) > 0 {
		userId := types.ZeroUid
		if types.GetTopicCat(topic) == types.TopicCatMe {
			userId = types.ParseUserId(topic)
			topic = ""
		}

		return s.adp.Files().LinkAttachments(topic, userId, msgId, fids)
	}

	return nil
}
