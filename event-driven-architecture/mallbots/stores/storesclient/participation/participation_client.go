// Code generated by go-swagger; DO NOT EDIT.

package participation

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new participation API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for participation API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DisableParticipation(params *DisableParticipationParams, opts ...ClientOption) (*DisableParticipationOK, error)

	EnableParticipation(params *EnableParticipationParams, opts ...ClientOption) (*EnableParticipationOK, error)

	GetParticipatingStores(params *GetParticipatingStoresParams, opts ...ClientOption) (*GetParticipatingStoresOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
DisableParticipation disables store service participation
*/
func (a *Client) DisableParticipation(params *DisableParticipationParams, opts ...ClientOption) (*DisableParticipationOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDisableParticipationParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "disableParticipation",
		Method:             "DELETE",
		PathPattern:        "/api/stores/{id}/participating",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DisableParticipationReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DisableParticipationOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DisableParticipationDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
EnableParticipation enables store service participation
*/
func (a *Client) EnableParticipation(params *EnableParticipationParams, opts ...ClientOption) (*EnableParticipationOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewEnableParticipationParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "enableParticipation",
		Method:             "PUT",
		PathPattern:        "/api/stores/{id}/participating",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &EnableParticipationReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*EnableParticipationOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*EnableParticipationDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetParticipatingStores gets a list of participating stores
*/
func (a *Client) GetParticipatingStores(params *GetParticipatingStoresParams, opts ...ClientOption) (*GetParticipatingStoresOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetParticipatingStoresParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getParticipatingStores",
		Method:             "GET",
		PathPattern:        "/api/stores/participating",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetParticipatingStoresReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetParticipatingStoresOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetParticipatingStoresDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}