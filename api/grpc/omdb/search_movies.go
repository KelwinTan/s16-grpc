package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/KelwinTan/s16-grpc/api/proto/v1/omdb"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func validateSearchMovieRequest(req *omdb.SearchMoviesRequest) error {
	if req.GetQuery() == "" || len(req.GetQuery()) < 1 {
		return status.Errorf(codes.InvalidArgument, "query should be more than 3 characters")
	}

	searchType := strings.ToLower(req.GetType())
	if searchType != "movie" && searchType != "series" && searchType != "episode" {
		return status.Errorf(codes.InvalidArgument, "type should be either movie or series or episode")
	}

	if req.GetPage() < 1 || req.GetPage() > 100 {
		return status.Errorf(codes.InvalidArgument, "page should be between 1-100")
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
