package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/KelwinTan/s16-grpc/api/proto/v1/omdb"
	"github.com/go-resty/resty/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	DEFAULT_OMDB_URL = "https://www.omdbapi.com/"
	API_KEY          = "faf7e5bb"
)

type Server struct {
	omdb.OMDBServiceServer
}

func validateGetMovieByIDRequest(req *omdb.GetMovieByIDRequest) {

}

func (s *Server) GetMovieByID(ctx context.Context, req *omdb.GetMovieByIDRequest) (*omdb.GetMovieByIDResponse, error) {
	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"apikey": API_KEY,
			"i":      req.Id,
		}).
		Get(DEFAULT_OMDB_URL)

	if err != nil {
		return nil, err
	}

	var movieResp omdb.GetMovieByIDResponse
	unmarshalOptions := protojson.UnmarshalOptions{
		AllowPartial:   true, // Allow unmarshaling even if some required fields are missing
		DiscardUnknown: true, // Discard unknown fields from JSON
	}

	if err := unmarshalOptions.Unmarshal(resp.Body(), &movieResp); err != nil {
		return nil, err
	}

	return &movieResp, nil
}

func validateSearchMovieRequest(req *omdb.SearchMoviesRequest) error {
	if req.Query == "" || len(req.Query) < 1 {
		return fmt.Errorf("query should be more than 3 characters (valid options from omdbapi.com)")
	}

	searchType := strings.ToLower(req.Type)
	if searchType != "movie" && searchType != "series" && searchType != "episode" {
		return fmt.Errorf("type should be either movie or series or episode (valid options from omdbapi.com)")
	}

	if req.Page < 1 || req.Page > 100 {
		return fmt.Errorf("page should be between 1-100 (valid options from omdbapi.com)")
	}

	return nil
}

func (s *Server) SearchMovies(ctx context.Context, req *omdb.SearchMoviesRequest) (*omdb.SearchMoviesResponse, error) {

	err := validateSearchMovieRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"apikey": API_KEY,
			"s":      req.Query,
			"type":   req.Type,
			"page":   fmt.Sprintf("%d", req.Page),
		}).
		Get(DEFAULT_OMDB_URL)

	if err != nil {
		return nil, err
	}

	var movieResp omdb.SearchMoviesResponse
	unmarshalOptions := protojson.UnmarshalOptions{
		AllowPartial:   true, // Allow unmarshaling even if some required fields are missing
		DiscardUnknown: true, // Discard unknown fields from JSON
	}

	if err := unmarshalOptions.Unmarshal(resp.Body(), &movieResp); err != nil {
		return nil, err
	}

	return &movieResp, nil
}
