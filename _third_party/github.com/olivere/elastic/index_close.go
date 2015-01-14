// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"bosun.org/_third_party/github.com/olivere/elastic/uritemplates"
)

// CloseIndexService closes an index.
// See documentation at http://www.elasticsearch.org/guide/en/elasticsearch/reference/1.4/indices-open-close.html.
type CloseIndexService struct {
	client            *Client
	debug             bool
	pretty            bool
	index             string
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
	timeout           string
	masterTimeout     string
}

// NewCloseIndexService creates a new CloseIndexService.
func NewCloseIndexService(client *Client) *CloseIndexService {
	return &CloseIndexService{client: client}
}

// Index is the name of the index.
func (s *CloseIndexService) Index(index string) *CloseIndexService {
	s.index = index
	return s
}

// Timeout is an explicit operation timeout.
func (s *CloseIndexService) Timeout(timeout string) *CloseIndexService {
	s.timeout = timeout
	return s
}

// MasterTimeout specifies the timeout for connection to master.
func (s *CloseIndexService) MasterTimeout(masterTimeout string) *CloseIndexService {
	s.masterTimeout = masterTimeout
	return s
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *CloseIndexService) IgnoreUnavailable(ignoreUnavailable bool) *CloseIndexService {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices. (This includes `_all` string or when no indices have been specified).
func (s *CloseIndexService) AllowNoIndices(allowNoIndices bool) *CloseIndexService {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both.
func (s *CloseIndexService) ExpandWildcards(expandWildcards string) *CloseIndexService {
	s.expandWildcards = expandWildcards
	return s
}

// buildURL builds the URL for the operation.
func (s *CloseIndexService) buildURL() (string, error) {
	// Build URL
	urls, err := uritemplates.Expand("/{index}/_close", map[string]string{
		"index": s.index,
	})
	if err != nil {
		return "", err
	}

	// Add query string parameters
	params := url.Values{}
	if s.allowNoIndices != nil {
		params.Set("allowNoIndices", fmt.Sprintf("%v", *s.allowNoIndices))
	}
	if s.expandWildcards != "" {
		params.Set("expandWildcards", s.expandWildcards)
	}
	if s.timeout != "" {
		params.Set("timeout", s.timeout)
	}
	if s.masterTimeout != "" {
		params.Set("masterTimeout", s.masterTimeout)
	}
	if s.ignoreUnavailable != nil {
		params.Set("ignoreUnavailable", fmt.Sprintf("%v", *s.ignoreUnavailable))
	}
	if len(params) > 0 {
		urls += "?" + params.Encode()
	}

	return urls, nil
}

// Validate checks if the operation is valid.
func (s *CloseIndexService) Validate() error {
	var invalid []string
	if s.index == "" {
		invalid = append(invalid, "Index")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *CloseIndexService) Do() (*CloseIndexResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	urls, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Setup HTTP request
	req, err := s.client.NewRequest("POST", urls)
	if err != nil {
		return nil, err
	}

	// Debug output?
	if s.debug {
		out, _ := httputil.DumpRequestOut((*http.Request)(req), true)
		log.Printf("%s\n", string(out))
	}

	// Get HTTP response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Debug output?
	if s.debug {
		out, _ := httputil.DumpResponse(res, true)
		log.Printf("%s\n", string(out))
	}

	// Return operation response
	resp := new(CloseIndexResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CloseIndexResponse is the response of CloseIndexService.Do.
type CloseIndexResponse struct {
	Acknowledged bool `json:"acknowledged"`
}
