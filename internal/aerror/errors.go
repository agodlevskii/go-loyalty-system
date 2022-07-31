package aerror

import (
	"fmt"
)

const (
	AccrualGet               = `unable to gather the order information from the accrual API`
	BalanceTableCreate       = `unable to create the balance table`
	BalanceSet               = `unable to set the user's balance`
	BalanceGet               = `unable to get the user's balance`
	BalanceInsufficient      = `the user balance is insufficient for the purchase`
	OrderTableCreate         = `unable to create the orders table`
	OrderAdd                 = `unable to create the order record`
	OrderNumberInvalid       = `the order number is invalid`
	OrderFind                = `unable to find the order record`
	OrderFindAll             = `unable to find the user order records`
	OrderUpdate              = `unable to update the order record`
	OrderExistsOtherUser     = `the order was added by another user`
	OrderExistsSameUser      = `the order is already enqueued by the user`
	RepoCreate               = `unable to initiate the application repo`
	WithdrawalTableCreate    = `unable to create the withdrawals table`
	WithdrawalAdd            = `unable to create the withdrawal record`
	WithdrawalFind           = `unable to find the withdrawal record`
	WithdrawalFindAll        = `unable to create the user withdrawal records`
	WithdrawalRequestInvalid = `the withdrawal request data is invalid`
)

type AppError struct {
	Label string
	Err   error
}

func NewError(label string, err error) *AppError {
	return &AppError{
		Label: label,
		Err:   err,
	}
}

func (e AppError) Error() string {
	if e.Error() == `` {
		return e.Label
	}
	return fmt.Sprintf("[%s] %v", e.Label, e.Error())
}

func (e AppError) Unwrap() error {
	return e.Err
}