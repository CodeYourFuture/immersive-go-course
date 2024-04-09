package main

import (
	"net/http"
	"net/http/httptest"
	"server-database/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchImages(t *testing.T) {
	expected := []types.Image{
		{
			Title:   "Sunset",
			AltText: "Clouds at sunset",
			URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
		{
			Title:   "Mountain",
			AltText: "A mountain at sunset",
			URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
	}
	actual, err := fetchImages()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestGetIndentParam(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "case 1:indent=0",
			query:    "indent=0",
			expected: "",
		},
		{
			name:     "case 2:indent=2",
			query:    "indent=2",
			expected: "  ",
		},
		{
			name:     "case 3:indent=4",
			query:    "indent=4",
			expected: "    ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/images.json?"+tc.query, nil)
			actual := getIndentParam(req)
			require.Equal(t, tc.expected, actual)
		})
	}
}
