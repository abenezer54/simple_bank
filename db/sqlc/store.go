package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()

}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// steps for transfer to succeed :
// create transfer
// create from account entry
// create to account entry
// update the balance of from and to accounts

func (store *Store) TransferTx(ctx context.Context, args CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	callback := func(q *Queries) error {
		transfer, err := q.CreateTransfer(context.Background(),
			CreateTransferParams{
				FromAccountID: args.FromAccountID,
				ToAccountID:   args.ToAccountID,
				Amount:        args.Amount,
			})
		if err != nil {
			return err
		}
		result.Transfer = transfer

		from_entry, err := q.CreateEntry(context.Background(),
			CreateEntryParams{
				AccountID: args.FromAccountID,
				Amount:    -args.Amount},
		)
		if err != nil {
			return err
		}
		result.FromEntry = from_entry

		to_entry, err := q.CreateEntry(context.Background(),
			CreateEntryParams{
				AccountID: args.ToAccountID,
				Amount:    args.Amount})

		if err != nil {
			return err
		}
		result.ToEntry = to_entry

		//  updating the balances
		account1, err := q.GetAccountForUpdate(context.Background(), args.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      account1.ID,
			Balance: account1.Balance - args.Amount,
		})
		if err != nil {
			return nil
		}

		account2, err := q.GetAccountForUpdate(context.Background(), args.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      account2.ID,
			Balance: account2.Balance + args.Amount,
		})
		if err != nil {
			return nil
		}
		return nil

	}

	//executing the callback
	err := store.execTx(context.Background(), callback)

	return result, err
}
