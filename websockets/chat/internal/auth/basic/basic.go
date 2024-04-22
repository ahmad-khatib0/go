package basic

import (
	"errors"
	"regexp"
	"strings"
	"time"

	t "github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	dt "github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"golang.org/x/crypto/bcrypt"
)

const realName = "basic"

// Define default constraints on login and password
const (
	defaultMinLoginLength    = 2
	defaultMaxLoginLength    = 32
	defaultMinPasswordLength = 3
)

// Token suitable as a login: starts and ends with a Unicode letter (class L) or number (class N),
//
// contains Unicode letters, numbers, dot and underscore.
var loginPattern = regexp.MustCompile(`^[\pL\pN][_.\pL\pN]*[\pL\pN]+$`)

// authenticator is the type to map authentication methods to.
type authenticator struct {
	db                dt.Adapter
	name              string
	cfg               config.AuthConfigBasic
	addToTags         bool
	minPasswordLength int
	minLoginLength    int
}

// Init initializes the basic authenticator.
func (a *authenticator) Init(configInterface interface{}, name string) error {
	cfg, ok := configInterface.(config.AuthConfigBasic)
	if !ok {
		return errors.New("basic_auth: failed to cast the config")
	}

	if name == "" {
		return errors.New("auth_basic: authenticator name cannot be blank")
	}

	if a.name != "" {
		return errors.New("auth_basic: already initialized as " + a.name + "; " + name)
	}

	a.name = name
	a.addToTags = cfg.AddToTags
	a.minPasswordLength = cfg.MinPasswordLength
	a.minLoginLength = cfg.MinLoginLength

	if a.minPasswordLength <= 0 {
		a.minPasswordLength = defaultMinPasswordLength
	}

	if a.minLoginLength > defaultMaxLoginLength {
		return errors.New("auth_basic: min_login_length exceeds the limit")
	}

	if a.minLoginLength <= 0 {
		a.minLoginLength = defaultMinLoginLength
	}

	return nil
}

// IsInitialized returns true if the handler is initialized.
func (a *authenticator) IsInitialized() bool {
	return a.name != ""
}

// AddRecord adds a basic authentication record to DB.
func (a *authenticator) AddRecord(rec *t.Rec, secret []byte, remoteAddr string) (*t.Rec, error) {
	uname, password, err := parseSecret(secret)
	if err != nil {
		return nil, err
	}

	if err = a.checkLoginPolicy(uname); err != nil {
		return nil, err
	}

	if err = a.checkPasswordPolicy(password); err != nil {
		return nil, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	var expires time.Time
	if rec.Lifetime > 0 {
		expires = time.Now().Add(time.Duration(rec.Lifetime)).UTC().Round(time.Millisecond)
	}

	authLevel := rec.AuthLevel
	if authLevel == t.LevelNone {
		authLevel = t.LevelAuth
	}

	err = a.db.Auth().AddRecord(rec.Uid, a.name, a.name+":"+uname, authLevel, passHash, expires)
	if err != nil {
		return nil, err
	}

	rec.AuthLevel = authLevel
	if a.addToTags {
		rec.Tags = append(rec.Tags, a.name+":"+uname)
	}

	return rec, nil
}

// UpdateRecord updates password for basic authentication.
func (a *authenticator) UpdateRecord(rec *t.Rec, secret []byte, remoteAddr string) (*t.Rec, error) {
	uname, password, err := parseSecret(secret)
	if err != nil {
		return nil, err
	}

	login, authLevel, _, _, err := a.db.Auth().GetRecord(rec.Uid, a.name)
	if err != nil {
		return nil, err
	}

	// User does not have a record.
	if login == "" {
		return nil, types.ErrNotFound
	}

	if uname == "" || uname == login {
		// User is changing just the password.
		uname = login
	} else if err = a.checkLoginPolicy(uname); err != nil {
		return nil, err

	} else if uid, _, _, _, err := a.db.Auth().GetUniqueRecord(a.name + ":" + uname); err != nil {
		return nil, err

	} else if !uid.IsZero() {
		// The (new) user name already exists. Report an error.
		return nil, types.ErrDuplicate
	}

	if err = a.checkPasswordPolicy(password); err != nil {
		return nil, err
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, types.ErrInternal
	}

	var expires time.Time
	if rec.Lifetime > 0 {
		expires = types.TimeNow().Add(time.Duration(rec.Lifetime))
	}

	err = a.db.Auth().UpdRecord(rec.Uid, a.name, a.name+":"+uname, authLevel, passhash, expires)
	if err != nil {
		return nil, err
	}

	// Remove old tag from the list of tags
	oldTag := a.name + ":" + login
	for i, tag := range rec.Tags {
		if tag == oldTag {
			rec.Tags[i] = rec.Tags[len(rec.Tags)-1]
			rec.Tags = rec.Tags[:len(rec.Tags)-1]
			break
		}
	}

	// Add new tag
	rec.Tags = append(rec.Tags, a.name+":"+uname)
	return rec, nil
}

// Authenticate checks login and password.
func (a *authenticator) Authenticate(secret []byte, remoteAddr string) (*t.Rec, []byte, error) {
	uname, password, err := parseSecret(secret)
	if err != nil {
		return nil, nil, err
	}

	uid, authLvl, passhash, expires, err := a.db.Auth().GetUniqueRecord(a.name + ":" + uname)
	if err != nil {
		return nil, nil, err
	}

	if uid.IsZero() {
		// Invalid login.
		return nil, nil, types.ErrFailed
	}

	if !expires.IsZero() && expires.Before(time.Now()) {
		// The record has expired
		return nil, nil, types.ErrExpired
	}

	err = bcrypt.CompareHashAndPassword(passhash, []byte(password))
	if err != nil {
		// Invalid password
		return nil, nil, types.ErrFailed
	}

	var lifetime time.Duration
	if !expires.IsZero() {
		lifetime = time.Until(expires)
	}

	return &t.Rec{
		Uid:       uid,
		AuthLevel: authLvl,
		Lifetime:  t.Duration(lifetime),
		Features:  0,
		State:     types.StateUndefined}, nil, nil
}

// AsTag convert search token into a prefixed tag, if possible.
func (a *authenticator) AsTag(token string) string {
	if !a.addToTags {
		return ""
	}

	if err := a.checkLoginPolicy(token); err != nil {
		return ""
	}

	return a.name + ":" + token
}

// IsUnique checks login uniqueness.
func (a *authenticator) IsUnique(secret []byte, remoteAddr string) (bool, error) {
	uname, _, err := parseSecret(secret)
	if err != nil {
		return false, err
	}

	if err := a.checkLoginPolicy(uname); err != nil {
		return false, err
	}

	uid, _, _, _, err := a.db.Auth().GetUniqueRecord(a.name + ":" + uname)
	if err != nil {
		return false, err
	}

	if uid.IsZero() {
		return true, nil
	}

	return false, types.ErrDuplicate
}

// GenSecret is not supported, generates an error.
func (authenticator) GenSecret(rec *t.Rec) ([]byte, time.Time, error) {
	return nil, time.Time{}, types.ErrUnsupported
}

// DelRecords deletes saved authentication records of the given user.
func (a *authenticator) DelRecords(uid types.Uid) error {
	return a.db.Auth().DelScheme(uid, a.name)
}

// RestrictedTags returns tag namespaces (prefixes) restricted by this adapter.
func (a *authenticator) RestrictedTags() ([]string, error) {
	var prefix []string
	if a.addToTags {
		prefix = []string{a.name}
	}

	return prefix, nil
}

// GetResetParams returns authenticator parameters passed to password reset handler.
func (a *authenticator) GetResetParams(uid types.Uid) (map[string]interface{}, error) {
	login, _, _, _, err := a.db.Auth().GetRecord(uid, a.name)

	if err != nil {
		return nil, err
	}

	// User does not have a record matching the authentication scheme.
	if login == "" {
		return nil, types.ErrNotFound
	}

	params := make(map[string]interface{})
	params["login"] = login

	return params, nil
}

// GetRealName returns the hardcoded name of the authenticator.
func (authenticator) GetRealName() string {
	return realName
}

// GetAuthConfig implements types.AuthHandler.
func (a *authenticator) GetAuthConfig() interface{} {
	return a.cfg
}

func (a *authenticator) checkLoginPolicy(uname string) error {
	rlogin := []rune(uname)
	if len(rlogin) < a.minLoginLength || len(rlogin) > defaultMaxLoginLength || !loginPattern.MatchString(uname) {
		return types.ErrPolicy
	}

	return nil
}

func (a *authenticator) checkPasswordPolicy(password string) error {
	if len([]rune(password)) < a.minPasswordLength {
		return types.ErrPolicy
	}

	return nil
}

// parseSecret() bsecret must be in the format "username:password"
func parseSecret(bsecret []byte) (uname, password string, err error) {
	secret := string(bsecret)

	splitAt := strings.Index(secret, ":")
	if splitAt < 0 {
		err = types.ErrMalformed
		return
	}

	uname = strings.ToLower(secret[:splitAt])
	password = secret[splitAt+1:]
	return
}
