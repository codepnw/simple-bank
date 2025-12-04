package transfergrpc

import (
	"context"

	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
	pb "github.com/codepnw/simple-bank/pb/proto"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransferServer struct {
	pb.UnimplementedSimpleBankServer
	uc transferusecase.TransferUsecase
}

func NewTransferServer(uc transferusecase.TransferUsecase) *TransferServer {
	return &TransferServer{uc: uc}
}

func (s *TransferServer) CreateTransfer(ctx context.Context, req *pb.CreateTransferRequest) (*pb.CreateTransferResponse, error) {
	input := &transferusecase.TransferParams{
		FromAccountID: req.GetFromAccountId(),
		ToAccountID:   req.GetToAccountId(),
		Amount:        req.GetAmount(),
		Currency:      req.GetCurrency(),
	}

	data, err := s.uc.Transfer(ctx, input)
	if err != nil {
		switch err {
		case errs.ErrAccountNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case errs.ErrTransferToSelf:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errs.ErrCurrencyMismatch:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errs.ErrMoneyNotEnough:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp := &pb.CreateTransferResponse{
		Transfer: &pb.Transfer{
			Id:            data.Transfer.ID,
			FromAccountId: data.Transfer.FromAccountID,
			ToAccountId:   data.Transfer.ToAccountID,
			Amount:        data.Transfer.Amount,
			CreatedAt:     timestamppb.New(data.Transfer.CreatedAt),
		},
		FromAccount: &pb.Account{
			Id:       data.FromAccount.ID,
			OwnerId:  data.FromAccount.OwnerID,
			Balance:  data.FromAccount.Balance,
			Currency: string(data.FromAccount.Currency),
		},
		ToAccount: &pb.Account{
			Id:       data.ToAccount.ID,
			OwnerId:  data.ToAccount.OwnerID,
			Balance:  data.ToAccount.Balance,
			Currency: string(data.ToAccount.Currency),
		},
		FromEntry: &pb.Entry{
			Id:        data.FromEntry.ID,
			AccountId: data.FromEntry.AccountID,
			Amount:    data.FromEntry.Amount,
			CreatedAt: timestamppb.New(data.Transfer.CreatedAt),
		},
		ToEntry: &pb.Entry{
			Id:        data.ToEntry.ID,
			AccountId: data.ToEntry.AccountID,
			Amount:    data.ToEntry.Amount,
			CreatedAt: timestamppb.New(data.Transfer.CreatedAt),
		},
	}
	return resp, nil
}
