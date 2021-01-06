package beaconv1

import (
	"context"
	"errors"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetForkSchedule retrieve all scheduled upcoming forks this node is aware of.
func (bs *Server) GetForkSchedule(ctx context.Context, req *emptypb.Empty) (*ethpb.ForkScheduleResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetSpec retrieves specification configuration (without Phase 1 params) used on this node. Specification params list
// Values are returned with following format:
// - any value starting with 0x in the spec is returned as a hex string.
// - all other values are returned as number.
func (bs *Server) GetSpec(ctx context.Context, req *emptypb.Empty) (*ethpb.SpecResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetDepositContract retrieves deposit contract address and genesis fork version.
func (bs *Server) GetDepositContract(ctx context.Context, req *emptypb.Empty) (*ethpb.DepositContractResponse, error) {
	return nil, errors.New("unimplemented")
}
