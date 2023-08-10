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
	"github.com/stretchr/testify/mock"
)

const (
	DEFAULT_TEST_OMDB_URL = "https://www.omdbapi.com/"
)

// MockMyService is a mock implementation of the MyService interface.
type MockOMDBService struct {
	mock.Mock
}

// GetMovieByID is a mocked implementation of the GetMovieByID gRPC method.
func (m *MockOMDBService) GetMovieByID(ctx context.Context, req *omdb.GetMovieByIDRequest) (*omdb.GetMovieByIDResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*omdb.GetMovieByIDResponse), args.Error(1)
}

// MockRestyClient is a mock implementation of the resty.Client interface.
type MockRestyClient struct {
	mock.Mock
}

func (m *MockRestyClient) R() *resty.Request {
	args := m.Called()
	return args.Get(0).(*resty.Request)
}

func (m *MockRestyClient) Get(url string) (*resty.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*resty.Response), args.Error(1)
}

// still need to mock gRPC service
func TestGetMovieByID_Success(t *testing.T) {
	mockClient := new(MockRestyClient)
	mockServer := new(MockOMDBService)

	req := &omdb.GetMovieByIDRequest{
		Id: "zzzzz",
	}

	expectedResponse := &omdb.GetMovieByIDResponse{
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

	// Set up the expected behavior of the mock.
	mockClient.On("R").Return(&resty.Request{})
	mockClient.On("Get", DEFAULT_TEST_OMDB_URL).Return(&resty.Response{
		RawResponse: &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(``))),
		},
	}, nil)

	mockServer.On("GetMovieByID", context.Background(), req).Return(expectedResponse, nil)

	actualResponse, err := mockServer.GetMovieByID(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
	mockClient.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestGetMovieByID_Fail_ID(t *testing.T) {
	req := &omdb.GetMovieByIDRequest{
		Id: "",
	}

	server := grpcS16.Server{}
	response, err := server.GetMovieByID(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = id should be filled and not be empty", err.Error())
	assert.Nil(t, response)
}
