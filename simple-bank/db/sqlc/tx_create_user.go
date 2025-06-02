package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
	// This function will be executed after the user is inserted, inside the same transaction.
	// And its output error will be used to decide whether to commit or rollback the transaction.
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		// the transaction will be rolled back if the an error occurred  here or above
		return arg.AfterCreate(result.User)
	})

	return result, err
}
