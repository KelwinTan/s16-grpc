package grpc

import (
	"github.com/KelwinTan/s16-grpc/api/proto/v1/omdb"
)

const (
	DEFAULT_OMDB_URL = "https://www.omdbapi.com/"
	API_KEY          = "faf7e5bb"
)

type Server struct {
	omdb.OMDBServiceServer
}
