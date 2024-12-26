package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/value"
	"github.com/samber/do"
	"time"
)

type Auth interface {
	GetRegistrationChallenge() (value.WebauthnRegistrationChallengeData, value.WebauthnRegistrationSessionData, error)
	Register(ctx context.Context, response value.WebauthnRegistrationChallengeResponseData, session value.WebauthnRegistrationSessionData) (entity.UserHandle, error)

	GetLoginChallenge() (value.WebauthnLoginChallengeData, value.WebauthnLoginSessionData, error)
	Login(ctx context.Context, response value.WebauthnLoginChallengeResponseData, session value.WebauthnLoginSessionData) (entity.UserHandle, error)
}

type authimpl struct {
	webauthn *webauthn.WebAuthn
	userRepo repository.User
}

func NewAuthService(injector *do.Injector) (Auth, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	loginTimeout, err := time.ParseDuration(config.Webauthn.Timeout.Login)
	if err != nil {
		return nil, fmt.Errorf("parsing webauthn login timeout: %w", err)
	}

	registrationTimeout, err := time.ParseDuration(config.Webauthn.Timeout.Registration)
	if err != nil {
		return nil, fmt.Errorf("parsing webauthn registration timeout: %w", err)
	}

	webauthn, err := webauthn.New(&webauthn.Config{
		RPID:                  config.Webauthn.RelyingParty.ID,
		RPDisplayName:         config.Webauthn.RelyingParty.DisplayName,
		RPOrigins:             config.Webauthn.RelyingParty.AllowedOrigins,
		AttestationPreference: protocol.PreferNoAttestation,
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    loginTimeout,
				TimeoutUVD: loginTimeout,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    registrationTimeout,
				TimeoutUVD: registrationTimeout,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating webauthn instance: %w", err)
	}

	userRepo, err := do.Invoke[repository.User](injector)
	if err != nil {
		return nil, fmt.Errorf("getting user repository: %w", err)
	}

	return &authimpl{
		webauthn: webauthn,
		userRepo: userRepo,
	}, nil
}

func (a *authimpl) GetRegistrationChallenge() (value.WebauthnRegistrationChallengeData, value.WebauthnRegistrationSessionData, error) {
	user, err := entity.NewUserWithRandomID()
	if err != nil {
		return nil, value.WebauthnRegistrationSessionData{}, fmt.Errorf("creating new anonymous user: %w", err)
	}

	options, webauthnSession, err := a.webauthn.BeginRegistration(user)
	if err != nil {
		return nil, value.WebauthnRegistrationSessionData{}, fmt.Errorf("starting usernameless registration: %w", err)
	}

	challenge, err := json.Marshal(options.Response)
	if err != nil {
		return nil, value.WebauthnRegistrationSessionData{}, fmt.Errorf("marshalling registration challenge: %w", err)
	}

	session, err := json.Marshal(webauthnSession)
	if err != nil {
		return nil, value.WebauthnRegistrationSessionData{}, fmt.Errorf("marshalling registration webauthnSession: %w", err)
	}

	return challenge, value.WebauthnRegistrationSessionData{
		Raw:    session,
		Expiry: webauthnSession.Expires,
	}, nil
}

func (a *authimpl) Register(ctx context.Context, response value.WebauthnRegistrationChallengeResponseData, session value.WebauthnRegistrationSessionData) (entity.UserHandle, error) {
	parsedResponse, err := protocol.ParseCredentialCreationResponseBytes(response)
	if err != nil {
		return nil, fmt.Errorf("parsing credential creation response: %w", err)
	}

	var webauthnSession webauthn.SessionData
	err = json.Unmarshal(session.Raw, &webauthnSession)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling webauthn session: %w", err)
	}

	tUser := entity.NewUser(webauthnSession.UserID, nil)

	credential, err := a.webauthn.CreateCredential(tUser, webauthnSession, parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("parsing credential: %w", err)
	}

	tUser.AddCredential(*credential)

	err = a.userRepo.AddUser(ctx, tUser)
	if err != nil {
		return nil, fmt.Errorf("inserting user: %w", err)
	}

	return tUser.Handle, nil
}

func (a *authimpl) GetLoginChallenge() (value.WebauthnLoginChallengeData, value.WebauthnLoginSessionData, error) {
	options, webauthnSession, err := a.webauthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, value.WebauthnLoginSessionData{}, fmt.Errorf("starting usernameless registration: %w", err)
	}

	challenge, err := json.Marshal(options.Response)
	if err != nil {
		return nil, value.WebauthnLoginSessionData{}, fmt.Errorf("marshalling registration challenge: %w", err)
	}

	session, err := json.Marshal(webauthnSession)
	if err != nil {
		return nil, value.WebauthnLoginSessionData{}, fmt.Errorf("marshalling registration webauthnSession: %w", err)
	}

	return challenge, value.WebauthnLoginSessionData{
		Raw:    session,
		Expiry: webauthnSession.Expires,
	}, nil
}

func (a *authimpl) Login(ctx context.Context, response value.WebauthnLoginChallengeResponseData, session value.WebauthnLoginSessionData) (entity.UserHandle, error) {
	parsedResponse, err := protocol.ParseCredentialRequestResponseBytes(response)
	if err != nil {
		return nil, fmt.Errorf("parsing credential request response: %w", err)
	}

	var webauthnSession webauthn.SessionData
	err = json.Unmarshal(session.Raw, &webauthnSession)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling webauthn session: %w", err)
	}

	var user *entity.User
	_, err = a.webauthn.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		user, err = a.userRepo.GetUser(ctx, userHandle)
		return user, err
	}, webauthnSession, parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("validating discoverable login: %w", err)
	}

	return user.Handle, nil
}
