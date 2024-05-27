package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// UpsertCred adds or updates a credential validation request. Return true if the record was inserted, false if updated.
func (s *Store) CredUpsert(cred *types.Credential) (bool, error) {
	cred.InitTimes()
	return s.adp.Credentials().Upsert(cred)
}

// ConfirmCred marks credential method as confirmed.
func (s *Store) CredConfirm(id types.Uid, method string) error {
	return s.adp.Credentials().Confirm(id, method)
}

// FailCred increments fail count for the given credential method.
func (s *Store) CredFail(id types.Uid, method string) error {
	return s.adp.Credentials().Fail(id, method)
}

// GetActiveCred gets a the currently active credential for the given user and method.
func (s *Store) CredGetActive(id types.Uid, method string) (*types.Credential, error) {
	return s.adp.Credentials().GetActive(id, method)
}

// GetAllCreds returns credentials of the given user, all or validated only.
func (s *Store) CredGetAllCreds(id types.Uid, method string, validatedOnly bool) ([]types.Credential, error) {
	return s.adp.Credentials().GetAll(id, method, validatedOnly)
}

// DelCred deletes user's credentials. If method is "", all credentials are deleted.
func (s *Store) CredDel(id types.Uid, method, value string) error {
	return s.adp.Credentials().Del(id, method, value)
}
