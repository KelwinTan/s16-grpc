package grpc_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	grpcS16 "github.com/KelwinTan/s16-grpc/api/grpc/omdb"
	"github.com/KelwinTan/s16-grpc/api/proto/v1/omdb"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// SearchMovies is a mocked implementation of the SearchMovies gRPC method.
func (m *MockOMDBService) SearchMovies(ctx context.Context, req *omdb.SearchMoviesRequest) (*omdb.SearchMoviesResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*omdb.SearchMoviesResponse), args.Error(1)
}

func TestSearchMovies_Fail_Query(t *testing.T) {
	req := &omdb.SearchMoviesRequest{
		Query: "",
		Type:  "",
		Page:  0,
	}

	server := grpcS16.Server{}
	response, err := server.SearchMovies(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = query should be more than 3 characters", err.Error())
	assert.Nil(t, response)
}

func TestSearchMovies_Fail_Type(t *testing.T) {
	req := &omdb.SearchMoviesRequest{
		Query: "Batman",
		Type:  "",
		Page:  0,
	}

	server := grpcS16.Server{}
	response, err := server.SearchMovies(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = type should be either movie or series or episode", err.Error())
	assert.Nil(t, response)
}

func TestSearchMovies_Fail_Page(t *testing.T) {
	req := &omdb.SearchMoviesRequest{
		Query: "Batman",
		Type:  "Movie",
		Page:  0,
	}

	server := grpcS16.Server{}
	response, err := server.SearchMovies(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = page should be between 1-100", err.Error())
	assert.Nil(t, response)
}

// still need to mock gRPC service
func TestSearchMovieByID_Success(t *testing.T) {
	mockClient := new(MockRestyClient)
	mockServer := new(MockOMDBService)

	req := &omdb.SearchMoviesRequest{
		Query: "",
	}

	expectedResponse := &omdb.SearchMoviesResponse{
		// Your expected response here
		// Id:        "tt4853102",
		// Title:     "Batman: The Killing Joke",
		// Year:      "2016",
		// Rated:     "R",
		// Genre:     "Animation, Action, Crime",
		// Plot:      "As Batman hunts for the escaped Joker, the Clown Prince of Crime attacks the Gordon family to prove a diabolical point mirroring his own fall into madness.",
		// Director:  "Sam Liu",
		// Actors:    []string(nil),
		// Language:  "English",
		// Country:   "United States",
		// Type:      "movie",
		// PosterUrl: "https://m.media-amazon.com/images/M/MV5BMTdjZTliODYtNWExMi00NjQ1LWIzN2MtN2Q5NTg5NTk3NzliL2ltYWdlXkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg",
	}

	mockClient.On("R").Return(&resty.Request{})
	mockClient.On("Get", DEFAULT_TEST_OMDB_URL).Return(&resty.Response{
		RawResponse: &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(``))),
		},
	}, nil)

	mockServer.On("SearchMovies", context.Background(), req).Return(expectedResponse, nil)

	actualResponse, err := mockServer.SearchMovies(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
	mockClient.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}
