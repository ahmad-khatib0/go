package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// Update updates a device record.
func (s *Store) DevUpdate(uid types.Uid, oldDeviceID string, dev *types.DeviceDef) error {
	// If the old device Id is specified and it's different from the new ID, delete the old id
	if oldDeviceID != "" && (dev == nil || dev.DeviceId != oldDeviceID) {
		if err := s.adp.Devices().Delete(uid, oldDeviceID); err != nil {
			return err
		}
	}

	// Insert or update the new DeviceId if one is given.
	if dev != nil && dev.DeviceId != "" {
		return s.adp.Devices().Upsert(uid, dev)
	}

	return nil
}

// GetAll returns all known device IDs for a given list of user IDs.
// The second return parameter is the count of found device IDs.
func (s *Store) DevGetAll(uid ...types.Uid) (map[types.Uid][]types.DeviceDef, int, error) {
	return s.adp.Devices().GetAll(uid...)
}

// Delete deletes device record for a given user.
func (s *Store) DevDelete(uid types.Uid, deviceID string) error {
	return s.adp.Devices().Delete(uid, deviceID)
}
