// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.
//
//go:build integration

package opensearchapi_test

import (
	"strings"
	"testing"

	"github.com/opensearch-project/opensearch-go/v3/opensearchapi"
	osapitest "github.com/opensearch-project/opensearch-go/v3/opensearchapi/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	client, err := opensearchapi.NewDefaultClient()
	require.Nil(t, err)

	index := "test-index-search"

	_, err = client.Index(
		nil,
		opensearchapi.IndexReq{
			Index:  index,
			Body:   strings.NewReader(`{"foo": "bar"}`),
			Params: opensearchapi.IndexParams{Refresh: "true"},
		},
	)
	require.Nil(t, err)
	t.Cleanup(func() {
		client.Indices.Delete(nil, opensearchapi.IndicesDeleteReq{Indices: []string{index}})
	})

	t.Run("with nil request", func(t *testing.T) {
		resp, err := client.Search(nil, nil)
		require.Nil(t, err)
		assert.NotNil(t, resp)
		osapitest.CompareRawJSONwithParsedJSON(t, resp, resp.Inspect().Response)
		assert.NotEmpty(t, resp.Hits.Hits)
	})

	t.Run("with request", func(t *testing.T) {
		resp, err := client.Search(nil, &opensearchapi.SearchReq{Indices: []string{index}, Body: strings.NewReader("")})
		require.Nil(t, err)
		assert.NotNil(t, resp)
		osapitest.CompareRawJSONwithParsedJSON(t, resp, resp.Inspect().Response)
		assert.NotEmpty(t, resp.Hits.Hits)
	})

	t.Run("inspect", func(t *testing.T) {
		failingClient, err := osapitest.CreateFailingClient()
		require.Nil(t, err)

		res, err := failingClient.Search(nil, nil)
		assert.NotNil(t, err)
		assert.NotNil(t, res)
		osapitest.VerifyInspect(t, res.Inspect())
	})
}
