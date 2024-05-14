package database

import (
	"context"
	"testing"

	util "github.com/fdhhhdjd/Banking_Platform_Golang/utils"
	"github.com/stretchr/testify/require"
)

func TestListAccounts(t *testing.T) {
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account,err := testQueries.CreateAccount(context.Background(),arg)
	
	require.NoError(t, err)
	require.NotEmpty(t,account)

	require.Equal(t, arg.Owner,account.Owner)
	require.Equal(t, arg.Balance,account.Balance)
	require.Equal(t, arg.Currency,account.Currency)

	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)

}