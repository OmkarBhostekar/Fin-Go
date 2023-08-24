package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "example.com/simplebank/db/sqlc"
	simplebank "example.com/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// check currency of both account
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}
	toAccount, valid := server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}
	if fromAccount.Currency != toAccount.Currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("currency mismatch between accounts")))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, ID int64, currency string) (account simplebank.Account, isvalid bool) {
	account, err := server.store.GetAcountById(ctx, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s vs %s", ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
