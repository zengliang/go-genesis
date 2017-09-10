// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package syspar

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"

	logger "github.com/EGaaS/go-egaas-mvp/packages/log"
	"github.com/EGaaS/go-egaas-mvp/packages/model"
)

const (
	// NumberNodes is the number of nodes
	NumberNodes = `number_of_dlt_nodes`
	// FuelRate is the rate
	FuelRate = `fuel_rate`
	// OpPrice is the costs of operations
	OpPrice = `op_price`
	// GapsBetweenBlocks is the time between blocks
	GapsBetweenBlocks = `gaps_between_blocks`
	// BlockchainURL is the address of the blockchain file.  For those who don't want to collect it from nodes
	BlockchainURL = `blockchain_url`
	// MaxBlockSize is the maximum size of the block
	MaxBlockSize = `max_block_size`
	// MaxTxSize is the maximum size of the transaction
	MaxTxSize = `max_tx_size`
	// MaxTxCount is the maximum count of the transactions
	MaxTxCount = `max_tx_count`
	// MaxColumns is the maximum columns in tables
	MaxColumns = `max_columns`
	// MaxIndexes is the maximum indexes in tables
	MaxIndexes = `max_indexes`
	// MaxBlockUserTx is the maximum number of user's transactions in one block
	MaxBlockUserTx = `max_block_user_tx`
	// UpdFullNodesPeriod is the maximum number of user's transactions in one block
	UpdFullNodesPeriod = `upd_full_nodes_period`
	// RecoveryAddress is the recovery address
	RecoveryAddress = `recovery_address`
	// CommissionWallet is the address for commissions
	CommissionWallet = `commission_wallet`
)

var (
	cache = map[string]string{
		BlockchainURL: "https://raw.githubusercontent.com/egaas-blockchain/egaas-blockchain.github.io/master/testnet_blockchain",
		// For compatible of develop versions
		// Remove later
		GapsBetweenBlocks:  `3`,
		MaxBlockSize:       `67108864`,
		MaxTxSize:          `33554432`,
		MaxTxCount:         `100000`,
		MaxColumns:         `50`,
		MaxIndexes:         `10`,
		MaxBlockUserTx:     `100`,
		UpdFullNodesPeriod: `3600`, // 3600 is for the test time, then we have to put 86400`
		RecoveryAddress:    `8275283526439353759`,
		CommissionWallet:   `8275283526439353759`,
	}
	cost  = make(map[string]int64)
	mutex = &sync.RWMutex{}
)

// SysUpdate reloads/updates values of system parameters
func SysUpdate() error {
	systemParameters, err := model.GetAllSystemParameters()
	if err != nil {
		logger.LogError(consts.DBError, err)
		return err
	}
	mutex.Lock()
	defer mutex.Unlock()
	for _, param := range systemParameters {
		cache[param.Name] = param.Value
	}

	cost = make(map[string]int64)
	json.Unmarshal([]byte(cache[OpPrice]), &cost)
	return err
}

func SysInt64(name string) int64 {
	val, err := strconv.ParseInt(SysString(name), 10, 64)
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(name))
	}
	return val
}

func GetBlockchainURL() string {
	return SysString(BlockchainURL)
}

func GetUpdFullNodesPeriod() int64 {
	val, err := strconv.ParseInt(SysString(UpdFullNodesPeriod), 10, 64)
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(UpdFullNodesPeriod))
	}
	return val
}

func GetMaxBlockSize() int64 {
	val, err := strconv.ParseInt(SysString(MaxBlockSize), 10, 64)
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(MaxBlockSize))
	}
	return val
}

func GetMaxTxSize() int64 {
	val, err := strconv.ParseInt(SysString(MaxTxSize), 10, 64)
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(MaxTxSize))
	}
	return val
}

func GetRecoveryAddress() int64 {
	val, err := strconv.ParseInt(SysString(RecoveryAddress), 10, 64)
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(RecoveryAddress))
	}
	return val
}

func GetCommissionWallet() int64 {
	val, err := strconv.ParseInt(SysString(CommissionWallet), 10, 64)
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(CommissionWallet))
	}
	return val
}

func GetGapsBetweenBlocks() int {
	val, err := strconv.Atoi(SysString(GapsBetweenBlocks))
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(GapsBetweenBlocks))
	}
	return val
}

func GetMaxTxCount() int {
	val, err := strconv.Atoi(SysString(MaxTxCount))
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(MaxTxCount))
	}
	return val
}

func GetMaxColumns() int {
	val, err := strconv.Atoi(SysString(MaxColumns))
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(MaxColumns))
	}
	return val
}

func GetMaxIndexes() int {
	val, err := strconv.Atoi(SysString(MaxIndexes))
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(MaxIndexes))
	}
	return val
}

func GetMaxBlockUserTx() int {
	val, err := strconv.Atoi(SysString(MaxBlockUserTx))
	if err != nil {
		logger.LogInfo(consts.StrToIntError, SysString(MaxBlockUserTx))
	}
	return val
}

// SysCost returns the cost of the transaction
func SysCost(name string) int64 {
	return cost[name]
}

// SysString returns string value of the system parameter
func SysString(name string) string {
	mutex.RLock()
	ret := cache[name]
	mutex.RUnlock()
	return ret
}
