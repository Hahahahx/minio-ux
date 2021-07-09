// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 7.12.0: DO NOT EDIT

package esapi

import (
	"context"
	"io"
	"net/http"
	"strings"
)

func newMLUpdateDataFrameAnalyticsFunc(t Transport) MLUpdateDataFrameAnalytics {
	return func(id string, body io.Reader, o ...func(*MLUpdateDataFrameAnalyticsRequest)) (*Response, error) {
		var r = MLUpdateDataFrameAnalyticsRequest{DocumentID: id, Body: body}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// MLUpdateDataFrameAnalytics - Updates certain properties of a data frame analytics job.
//
// This API is beta.
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/current/update-dfanalytics.html.
//
type MLUpdateDataFrameAnalytics func(id string, body io.Reader, o ...func(*MLUpdateDataFrameAnalyticsRequest)) (*Response, error)

// MLUpdateDataFrameAnalyticsRequest configures the ML Update Data Frame Analytics API request.
//
type MLUpdateDataFrameAnalyticsRequest struct {
	DocumentID string

	Body io.Reader

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r MLUpdateDataFrameAnalyticsRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "POST"

	path.Grow(1 + len("_ml") + 1 + len("data_frame") + 1 + len("analytics") + 1 + len(r.DocumentID) + 1 + len("_update"))
	path.WriteString("/")
	path.WriteString("_ml")
	path.WriteString("/")
	path.WriteString("data_frame")
	path.WriteString("/")
	path.WriteString("analytics")
	path.WriteString("/")
	path.WriteString(r.DocumentID)
	path.WriteString("/")
	path.WriteString("_update")

	params = make(map[string]string)

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), r.Body)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if r.Body != nil {
		req.Header[headerContentType] = headerContentTypeJSON
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f MLUpdateDataFrameAnalytics) WithContext(v context.Context) func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f MLUpdateDataFrameAnalytics) WithPretty() func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f MLUpdateDataFrameAnalytics) WithHuman() func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f MLUpdateDataFrameAnalytics) WithErrorTrace() func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f MLUpdateDataFrameAnalytics) WithFilterPath(v ...string) func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f MLUpdateDataFrameAnalytics) WithHeader(h map[string]string) func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
//
func (f MLUpdateDataFrameAnalytics) WithOpaqueID(s string) func(*MLUpdateDataFrameAnalyticsRequest) {
	return func(r *MLUpdateDataFrameAnalyticsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
