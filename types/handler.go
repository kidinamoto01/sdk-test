package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"strconv"
//	"github.com/tendermint/abci/types"
)

// TODO: Admin cannot do ops other than create account msg

func RegisterBallotRoutes(r baseapp.Router, accts sdk.AccountMapper,storeKey sdk.StoreKey){
	r.AddRoute(ProposeVoteType, ProposeMsgHandler(accts))
	RegisterRoutes(r,accts)
}

//Business logic is executed here

// Routes the message (request) to a proper handler
func RegisterRoutes(r baseapp.Router, accts sdk.AccountMapper) {
	r.AddRoute(DepositType, DepositMsgHandler(accts))
	r.AddRoute(SettlementType, SettleMsgHandler(accts))
	r.AddRoute(WithdrawType, WithdrawMsgHandler(accts))
	r.AddRoute(CreateAccountType, CreateAccountMsgHandler(accts))
}

/*

Deposit functionality.

Sender is Custodian
Rec is Member
*/
func DepositMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return depositMsgHandler{accts}.Do
}

type depositMsgHandler struct {
	accts sdk.AccountMapper
}

// Deposit logic
func (d depositMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// TODO: ensure auth actually checks the sigs

	// ensure proper message
	dm, ok := msg.(DepositMsg)
	if !ok {
		return ErrWrongMsgFormat("expected DepositMsg").Result()
	}

	// ensure proper types
	sender, err := getAccountWithType(ctx, d.accts, dm.Sender, IsCustodian)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAccountWithType(ctx, d.accts, dm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}

	err = moveMoney(d.accts, ctx, sender, rcpt, dm.Amount, false, true)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

/*

Propose functionality.

*/
func ProposeMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return proposeMsgHandler{accts}.Do
}

type proposeMsgHandler struct {
	accts sdk.AccountMapper
}

// propose logic
func (d proposeMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {

	// ensure proper message
	dm, ok := msg.(ProposeMsg)
	if !ok {
		return ErrWrongMsgFormat("expected DepositMsg").Result()
	}

	// ensure proper types
	sender, err := getAccount(ctx, d.accts, dm.Sender)

	if err != nil {
		return err.Result()
	}
	key := sender.Address.Bytes()
	value := []byte(strconv.Itoa(dm.Index))

	fmt.Println("key: ",key," value: ",value)
	//skey := sdk.NewKVStoreKey("new")
	//store := ctx.KVStore(skey)
	//store.Set(key, value)

	return sdk.Result{Code:sdk.CodeOK}
}

//func (*sapp SearchApp)insertValue(ctx sdk.Context){
//
//
//}
/*
Settlement funcionality.

Sender is CH
Rec is member
*/
func SettleMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return settleMsgHandler{accts}.Do
}

type settleMsgHandler struct {
	accts sdk.AccountMapper
}

// Settlement logic
func (sh settleMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {

	// ensure proper message
	sm, ok := msg.(SettleMsg)
	if !ok {
		return ErrWrongMsgFormat("expected SettleMsg").Result()
	}

	// ensure proper types
	sender, err := getAccountWithType(ctx, sh.accts, sm.Sender, IsClearingHouse)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAccountWithType(ctx, sh.accts, sm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}

	err = moveMoney(sh.accts, ctx, sender, rcpt, sm.Amount, false, true)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}

}

/*

Withdraw functionality.

Sender is member
Reci is custodian
Operator is CH

*/
func WithdrawMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return withdrawMsgHandler{accts}.Do
}

type withdrawMsgHandler struct {
	accts sdk.AccountMapper
}

// Withdraw logic
func (wh withdrawMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	wm, ok := msg.(WithdrawMsg)
	if !ok {
		return ErrWrongMsgFormat("expected WithdrawMsg").Result()
	}

	// ensure proper types
	sender, err := getAccountWithType(ctx, wh.accts, wm.Sender, IsMember)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAccountWithType(ctx, wh.accts, wm.Recipient, IsCustodian)
	if err != nil {
		return err.Result()
	}
	_, err = getAccountWithType(ctx, wh.accts, wm.Operator, IsClearingHouse)
	if err != nil {
		return err.Result()
	}

	err = moveMoney(wh.accts, ctx, sender, rcpt, wm.Amount, true, false)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}

}

/*
Create account functionality.

Creator is CH
*/
func CreateAccountMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return createAccountMsgHandler{accts}.Do
}

type createAccountMsgHandler struct {
	accts sdk.AccountMapper
}

// Create acc logic.
// A clearing house account is allowed to create any kind of accounts,
// including clearing house, custodian, and members accounts.
// TODO: clarify business rules and ensure this is desired
func (h createAccountMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(CreateAccountMsg)
	if !ok {
		return ErrWrongMsgFormat("expected CreateAccountMsg").Result()
	}
	// ensure creator exists
	creator := h.accts.GetAccount(ctx, cm.Creator)
	if creator == nil {
		return ErrInvalidAccount("couldn't find creator").Result()
	}
	// ensure proper types
	concreteCreator := creator.(*AppAccount)
	// ensure new account does not exist
	if rawAccount := h.accts.GetAccount(ctx, cm.PubKey.Address()); rawAccount != nil {
		return ErrInvalidAccount("the account already exists").Result()
	}

	// Construct a new account
	acct := NewAppAccount(cm.PubKey, nil, cm.AccountType, creator.GetAddress(), cm.IsAdmin, cm.LegalEntityName)
	if err := CanCreate(concreteCreator, acct); err != nil {
		return ErrUnauthorized(fmt.Sprintf("can't create account: %v", err)).Result()
	}
	h.accts.SetAccount(ctx, acct)
	return sdk.Result{}
}

//*********************************** helper methods *********************************************

// Transfers money from the sender to the  recipient
func moveMoney(accts sdk.AccountMapper, ctx sdk.Context, sender *AppAccount, recipient *AppAccount,
	amount sdk.Coin, senderMustBePositive bool, recipientMustBePositive bool) sdk.Error {

	transfer := sdk.Coins{amount}

	// first verify funds
	sender.Coins = sender.Coins.Minus(transfer)
	if senderMustBePositive && !sender.Coins.IsNotNegative() {
		return ErrInvalidAmount("sender  has insufficient funds")
	}
	// transfer may be negative
	recipient.Coins = recipient.Coins.Plus(transfer)
	if recipientMustBePositive && !recipient.Coins.IsNotNegative() {
		return ErrInvalidAmount("recipient has insufficient funds")
	}

	// now make the transfer and save the result
	accts.SetAccount(ctx, sender)
	accts.SetAccount(ctx, recipient)
	return nil
}

// Returns the account and verifies its type
func getAccountWithType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	typeCheck func(*AppAccount) bool) (*AppAccount, sdk.Error) {

	rawAccount := accts.GetAccount(ctx, addr)
	if rawAccount == nil {
		return nil, ErrInvalidAccount("account does not exist")
	}
	account := rawAccount.(*AppAccount)
	if !typeCheck(account) {
		return nil, ErrWrongSigner(account.Type)
	}

	return account, nil
}


// Returns the account
func getAccount(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (
	*AppAccount, sdk.Error) {

	rawAccount := accts.GetAccount(ctx, addr)
	if rawAccount == nil {
		return nil, ErrInvalidAccount("account does not exist")
	}
	account := rawAccount.(*AppAccount)

	return account, nil
}

// Creates an account instance
func createAccount(creator crypto.Address, newAccPubKey crypto.PubKey, typ string, entity string, isAdmin bool) *AppAccount {
	acct := new(AppAccount)
	acct.SetAddress(newAccPubKey.Address())
	acct.SetPubKey(newAccPubKey)
	acct.SetCoins(nil)
	acct.SetCreator(creator)
	acct.Type = typ
	acct.LegalEntityName = entity
	acct.EntityAdmin = isAdmin
	return acct
}
