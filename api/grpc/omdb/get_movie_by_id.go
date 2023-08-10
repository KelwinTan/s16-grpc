package grpc

import (
	"context"

	"github.com/KelwinTan/s16-grpc/api/proto/v1/omdb"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func validateGetMovieByIDRequest(req *omdb.GetMovieByIDRequest) error {
	if req.GetId() == "" || len(req.GetId()) < 1 {
		return status.Errorf(codes.InvalidArgument, "id should be filled and not be empty")
	}

	return nil
}

func (s *Server) GetMovieByID(ctx context.Context, req *omdb.GetMovieByIDRequest) (*omdb.GetMovieByIDResponse, error) {
	err := validateGetMovieByIDRequest(req)
	if err != nil {
		return nil, err
	}

	var movieResp omdb.GetMovieByIDResponse
	unmarshalOptions := protojson.UnmarshalOptions{
		AllowPartial:   true, // Allow unmarshaling even if some required fields are missing
		DiscardUnknown: true, // Discard unknown fields from JSON
	}

	// tried caching using memcache but time is not enough

	// it, err := cache.MemCache.Get(req.Id)
	// if err != nil {
	// 	return nil, err
	// }

	// if it != nil {
	// 	if err := unmarshalOptions.Unmarshal(it.Value, &movieResp); err != nil {
	// 		return nil, err
	// 	}
	// 	return &movieResp, nil
	// }

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"apikey": API_KEY,
			"i":      req.Id,
		}).
		Get(DEFAULT_OMDB_URL)

	if err != nil {
		return nil, err
	}

	if err := unmarshalOptions.Unmarshal(resp.Body(), &movieResp); err != nil {
		return nil, err
	}

	// cache.MemCache.Set(&memcache.Item{
	// 	Key:   req.Id,
	// 	Value: resp.Body(),
	// })

	return &movieResp, nil
}
