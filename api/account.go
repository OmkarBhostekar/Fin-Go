package api

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"

	db "example.com/simplebank/db/sqlc"
	"example.com/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR INR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pkErr, ok := err.(*pq.Error); ok {
			switch pkErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := server.store.GetAcountById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account.Owner != authPayload.Username {
		err := fmt.Errorf("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type ListAccountResponse struct {
	Accounts []db.Account `json:"accounts"`
	PageID   int32        `json:"page_id"`
	Pages    int32        `json:"total_pages"`
}

func (server *Server) ListAccounts(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetAllAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Owner:  authPayload.Username,
	}
	accounts, err := server.store.GetAllAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if accounts == nil {
		accounts = []db.Account{}
	}
	totalAccounts, err := server.store.GetAccountsCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	pages := int32(math.Ceil(float64(totalAccounts) / float64(req.PageSize)))
	ctx.JSON(http.StatusOK, ListAccountResponse{
		Accounts: accounts,
		PageID:   req.PageID,
		Pages:    pages,
	})
}
