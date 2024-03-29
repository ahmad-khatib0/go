// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// StorespbGetStoresResponse storespb get stores response
//
// swagger:model storespbGetStoresResponse
type StorespbGetStoresResponse struct {

	// stores
	Stores []*StorespbStore `json:"stores"`
}

// Validate validates this storespb get stores response
func (m *StorespbGetStoresResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStores(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorespbGetStoresResponse) validateStores(formats strfmt.Registry) error {
	if swag.IsZero(m.Stores) { // not required
		return nil
	}

	for i := 0; i < len(m.Stores); i++ {
		if swag.IsZero(m.Stores[i]) { // not required
			continue
		}

		if m.Stores[i] != nil {
			if err := m.Stores[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("stores" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("stores" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this storespb get stores response based on the context it is used
func (m *StorespbGetStoresResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateStores(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorespbGetStoresResponse) contextValidateStores(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Stores); i++ {

		if m.Stores[i] != nil {

			if swag.IsZero(m.Stores[i]) { // not required
				return nil
			}

			if err := m.Stores[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("stores" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("stores" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *StorespbGetStoresResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorespbGetStoresResponse) UnmarshalBinary(b []byte) error {
	var res StorespbGetStoresResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
