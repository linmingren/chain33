package types

import (
	"errors"
)

var (
	ErrNotFound                   = errors.New("ErrNotFound")
	ErrBlockExec                  = errors.New("ErrBlockExec")
	ErrCheckStateHash             = errors.New("ErrCheckStateHash")
	ErrCheckTxHash                = errors.New("ErrCheckTxHash")
	ErrReRunGenesis               = errors.New("ErrReRunGenesis")
	ErrActionNotSupport           = errors.New("ErrActionNotSupport")
	ErrChannelFull                = errors.New("ErrChannelFull")
	ErrAmount                     = errors.New("ErrAmount")
	ErrNoTicket                   = errors.New("ErrNoTicket")
	ErrMinerIsStared              = errors.New("ErrMinerIsStared")
	ErrMinerNotStared             = errors.New("ErrMinerNotStared")
	ErrMinerNotClosed             = errors.New("ErrMinerNotClosed")
	ErrTicketCount                = errors.New("ErrTicketCount")
	ErrHashlockAmount             = errors.New("ErrHashlockAmount")
	ErrHashlockHash               = errors.New("ErrHashlockHash")
	ErrHashlockStatus             = errors.New("ErrHashlockStatus")
	ErrNoPeer                     = errors.New("ErrNoPeer")
	ErrExecNameNotMath            = errors.New("ErrExecNameNotMath")
	ErrChannelClosed              = errors.New("ErrChannelClosed")
	ErrNotMinered                 = errors.New("ErrNotMinered")
	ErrTime                       = errors.New("ErrTime")
	ErrFromAddr                   = errors.New("ErrFromAddr")
	ErrBlockHeight                = errors.New("ErrBlockHeight")
	ErrCoinBaseExecer             = errors.New("ErrCoinBaseExecer")
	ErrCoinBaseTxType             = errors.New("ErrCoinBaseTxType")
	ErrCoinBaseExecErr            = errors.New("ErrCoinBaseExecErr")
	ErrCoinBaseTarget             = errors.New("ErrCoinBaseTarget")
	ErrCoinbaseReward             = errors.New("ErrCoinbaseReward")
	ErrNotAllowDeposit            = errors.New("ErrNotAllowDeposit")
	ErrCoinBaseIndex              = errors.New("ErrCoinBaseIndex")
	ErrCoinBaseTicketStatus       = errors.New("ErrCoinBaseTicketStatus")
	ErrBlockNotFound              = errors.New("ErrBlockNotFound")
	ErrHashlockReturnAddrss       = errors.New("ErrHashlockReturnAddrss")
	ErrHashlockTime               = errors.New("ErrHashlockTime")
	ErrHashlockReapeathash        = errors.New("ErrHashlockReapeathash")
	ErrHashlockSendAddress        = errors.New("ErrHashlockSendAddress")
	ErrRetrieveRepeatAddress      = errors.New("ErrRetrieveRepeatAddress")
	ErrRetrievePeriodLimit        = errors.New("ErrRetrievePeriodLimit")
	ErrRetrieveAmountLimit        = errors.New("ErrRetrieveAmountLimit")
	ErrRetrieveTimeweightLimit    = errors.New("ErrRetrieveTimeweightLimit")
	ErrRetrievePrepareAddress     = errors.New("ErrRetrievePrepareAddress")
	ErrRetrievePerformAddress     = errors.New("ErrRetrievePerformAddress")
	ErrRetrieveCancelAddress      = errors.New("ErrRetrieveCancelAddress")
	ErrRetrieveStatus             = errors.New("ErrRetrieveStatus")
	ErrRetrieveRelateLimit        = errors.New("ErrRetrieveRelateLimit")
	ErrRetrieveRelation           = errors.New("ErrRetrieveRelation")
	ErrRetrieveNoBalance          = errors.New("ErrRetrieveNoBalance")
	ErrStartBigThanEnd            = errors.New("ErrStartBigThanEnd")
	ErrToAddrNotSameToExecAddr    = errors.New("ErrToAddrNotSameToExecAddr")
	ErrTypeAsset                  = errors.New("ErrTypeAsset")
	ErrEmpty                      = errors.New("ErrEmpty")
	ErrSendSameToRecv             = errors.New("ErrSendSameToRecv")
	ErrExecNameNotAllow           = errors.New("ErrExecNameNotAllow")
	ErrLocalDBPerfix              = errors.New("ErrLocalDBPerfix")
	ErrTimeout                    = errors.New("ErrTimeout")
	ErrBlockHeaderDifficulty      = errors.New("ErrBlockHeaderDifficulty")
	ErrNoTx                       = errors.New("ErrNoTx")
	ErrTxExist                    = errors.New("ErrTxExist")
	ErrManyTx                     = errors.New("ErrManyTx")
	ErrDupTx                      = errors.New("ErrDupTx")
	ErrMemFull                    = errors.New("ErrMemFull")
	ErrNoBalance                  = errors.New("ErrNoBalance")
	ErrBalanceLessThanTenTimesFee = errors.New("ErrBalanceLessThanTenTimesFee")
	ErrTxExpire                   = errors.New("ErrTxExpire")
	ErrSign                       = errors.New("ErrSign")
	ErrFeeTooLow                  = errors.New("ErrFeeTooLow")
	ErrEmptyTx                    = errors.New("ErrEmptyTx")
	ErrTxFeeTooLow                = errors.New("ErrTxFeeTooLow")
	ErrTxMsgSizeTooBig            = errors.New("ErrTxMsgSizeTooBig")
	ErrTicketClosed               = errors.New("ErrTicketClosed")
	ErrEmptyMinerTx               = errors.New("ErrEmptyMinerTx")
	ErrMinerNotPermit             = errors.New("ErrMinerNotPermit")
	ErrMinerAddr                  = errors.New("ErrMinerAddr")
	ErrModify                     = errors.New("ErrModify")
	ErrFutureBlock                = errors.New("ErrFutureBlock")
	ErrHashNotFound               = errors.New("ErrHashNotFound")
	ErrTxDup                      = errors.New("ErrTxDup")

	// BlockChain Error Types
	ErrHashNotExist           = errors.New("ErrHashNotExist")
	ErrHeightNotExist         = errors.New("ErrHeightNotExist")
	ErrTxNotExist             = errors.New("ErrTxNotExist")
	ErrAddrNotExist           = errors.New("ErrAddrNotExist")
	ErrStartHeight            = errors.New("ErrStartHeight")
	ErrEndLessThanStartHeight = errors.New("ErrEndLessThanStartHeight")
	ErrClientNotBindQueue     = errors.New("ErrClientNotBindQueue")
	ErrContinueBack           = errors.New("ErrContinueBack")
	ErrUnmarshal              = errors.New("ErrUnmarshal")
	ErrMarshal                = errors.New("ErrMarshal")
	ErrBlockExist             = errors.New("ErrBlockExist")
	ErrParentBlockNoExist     = errors.New("ErrParentBlockNoExist")
	ErrBlockHeightNoMatch     = errors.New("ErrBlockHeightNoEqual")
	ErrParentTdNoExist        = errors.New("ErrParentTdNoExist")
	ErrBlockHashNoMatch       = errors.New("ErrBlockHashNoMatch")
	ErrIsClosed               = errors.New("ErrIsClosed")
	ErrDecode                 = errors.New("ErrDecode")
	ErrNotRollBack            = errors.New("ErrNotRollBack")

	//wallet
	ErrInputPara      = errors.New("ErrInputPara")
	ErrWalletIsLocked = errors.New("ErrWalletIsLocked")
	ErrSaveSeedFirst  = errors.New("ErrSaveSeedFirst")
	ErrUnLockFirst    = errors.New("ErrUnLockFirst")

	ErrLabelHasUsed        = errors.New("ErrLabelHasUsed")
	ErrPrivkeyExist        = errors.New("ErrPrivkeyExist")
	ErrPrivkey             = errors.New("ErrPrivkey")
	ErrInsufficientBalance = errors.New("ErrInsufficientBalance")
	ErrVerifyOldpasswdFail = errors.New("ErrVerifyOldpasswdFail")
	ErrInputPassword       = errors.New("ErrInputPassword")
	ErrSeedlang            = errors.New("ErrSeedlang")
	ErrSeedNotExist        = errors.New("ErrSeedNotExist")
	ErrSubPubKeyVerifyFail = errors.New("ErrSubPubKeyVerifyFail")
	ErrLabelNotExist       = errors.New("ErrLabelNotExist")
	ErrAccountNotExist     = errors.New("ErrAccountNotExist")
	ErrSeedExist           = errors.New("ErrSeedExist")
	ErrNotSupport          = errors.New("ErrNotSupport")
	ErrSeedWordNum         = errors.New("ErrSeedWordNum")
	ErrOnlyTicketUnLocked  = errors.New("ErrOnlyTicketUnLocked")
	ErrNewCrypto           = errors.New("ErrNewCrypto")
	ErrFromHex             = errors.New("ErrFromHex")
	ErrPrivKeyFromBytes    = errors.New("ErrFromHex")
)