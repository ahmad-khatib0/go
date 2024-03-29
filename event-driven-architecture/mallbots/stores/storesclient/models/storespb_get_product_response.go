// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// StorespbGetProductResponse storespb get product response
//
// swagger:model storespbGetProductResponse
type StorespbGetProductResponse struct {

	// product
	Product *StorespbProduct `json:"product,omitempty"`
}

// Validate validates this storespb get product response
func (m *StorespbGetProductResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateProduct(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorespbGetProductResponse) validateProduct(formats strfmt.Registry) error {
	if swag.IsZero(m.Product) { // not required
		return nil
	}

	if m.Product != nil {
		if err := m.Product.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("product")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("product")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this storespb get product response based on the context it is used
func (m *StorespbGetProductResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateProduct(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorespbGetProductResponse) contextValidateProduct(ctx context.Context, formats strfmt.Registry) error {

	if m.Product != nil {

		if swag.IsZero(m.Product) { // not required
			return nil
		}

		if err := m.Product.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("product")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("product")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *StorespbGetProductResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorespbGetProductResponse) UnmarshalBinary(b []byte) error {
	var res StorespbGetProductResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
