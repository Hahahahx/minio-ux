// Code generated by go-swagger; DO NOT EDIT.

// This file is part of MinIO Console Server
// Copyright (c) 2021 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

package user_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/minio/console/models"
)

// LogoutOKCode is the HTTP code returned for type LogoutOK
const LogoutOKCode int = 200

/*LogoutOK A successful response.

swagger:response logoutOK
*/
type LogoutOK struct {
}

// NewLogoutOK creates LogoutOK with default headers values
func NewLogoutOK() *LogoutOK {

	return &LogoutOK{}
}

// WriteResponse to the client
func (o *LogoutOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

/*LogoutDefault Generic error response.

swagger:response logoutDefault
*/
type LogoutDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewLogoutDefault creates LogoutDefault with default headers values
func NewLogoutDefault(code int) *LogoutDefault {
	if code <= 0 {
		code = 500
	}

	return &LogoutDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the logout default response
func (o *LogoutDefault) WithStatusCode(code int) *LogoutDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the logout default response
func (o *LogoutDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the logout default response
func (o *LogoutDefault) WithPayload(payload *models.Error) *LogoutDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the logout default response
func (o *LogoutDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *LogoutDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
