// Code generated by go-swagger; DO NOT EDIT.

package shopping_list

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewCcompleteShoppingListParams creates a new CcompleteShoppingListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCcompleteShoppingListParams() *CcompleteShoppingListParams {
	return &CcompleteShoppingListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCcompleteShoppingListParamsWithTimeout creates a new CcompleteShoppingListParams object
// with the ability to set a timeout on a request.
func NewCcompleteShoppingListParamsWithTimeout(timeout time.Duration) *CcompleteShoppingListParams {
	return &CcompleteShoppingListParams{
		timeout: timeout,
	}
}

// NewCcompleteShoppingListParamsWithContext creates a new CcompleteShoppingListParams object
// with the ability to set a context for a request.
func NewCcompleteShoppingListParamsWithContext(ctx context.Context) *CcompleteShoppingListParams {
	return &CcompleteShoppingListParams{
		Context: ctx,
	}
}

// NewCcompleteShoppingListParamsWithHTTPClient creates a new CcompleteShoppingListParams object
// with the ability to set a custom HTTPClient for a request.
func NewCcompleteShoppingListParamsWithHTTPClient(client *http.Client) *CcompleteShoppingListParams {
	return &CcompleteShoppingListParams{
		HTTPClient: client,
	}
}

/*
CcompleteShoppingListParams contains all the parameters to send to the API endpoint

	for the ccomplete shopping list operation.

	Typically these are written to a http.Request.
*/
type CcompleteShoppingListParams struct {

	// Body.
	Body interface{}

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the ccomplete shopping list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CcompleteShoppingListParams) WithDefaults() *CcompleteShoppingListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the ccomplete shopping list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CcompleteShoppingListParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) WithTimeout(timeout time.Duration) *CcompleteShoppingListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) WithContext(ctx context.Context) *CcompleteShoppingListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) WithHTTPClient(client *http.Client) *CcompleteShoppingListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) WithBody(body interface{}) *CcompleteShoppingListParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) SetBody(body interface{}) {
	o.Body = body
}

// WithID adds the id to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) WithID(id string) *CcompleteShoppingListParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the ccomplete shopping list params
func (o *CcompleteShoppingListParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *CcompleteShoppingListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
