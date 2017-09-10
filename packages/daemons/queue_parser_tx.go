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

package daemons

import (
	"context"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	logger "github.com/EGaaS/go-egaas-mvp/packages/log"
	"github.com/EGaaS/go-egaas-mvp/packages/logging"
	"github.com/EGaaS/go-egaas-mvp/packages/model"
	"github.com/EGaaS/go-egaas-mvp/packages/parser"
)

// QueueParserTx parses transaction from the queue
func QueueParserTx(d *daemon, ctx context.Context) error {
	logger.LogDebug(consts.FuncStarted, "")
	lock, err := DbLock(ctx, d.goRoutineName)
	if !lock || err != nil {
		logger.LogError(consts.DBError, err)
		return err
	}
	defer DbUnlock(d.goRoutineName)

	infoBlock := &model.InfoBlock{}
	err = infoBlock.GetInfoBlock()
	if err != nil {
		logger.LogError(consts.DBError, err)
		return err
	}
	if infoBlock.BlockID == 0 {
		logger.LogDebug(consts.DebugMessage, "there are now blocks for parse")
		return nil
	}

	// delete looped transactions
	logging.WriteSelectiveLog("DELETE FROM transactions WHERE verified = 0 AND used = 0 AND counter > 10")
	affect, err := model.DeleteLoopedTransactions()
	if err != nil {
		logging.WriteSelectiveLog(err)
		logger.LogError(consts.DBError, err)
		return err
	}
	logging.WriteSelectiveLog("affect: " + converter.Int64ToStr(affect))

	p := new(parser.Parser)
	err = p.AllTxParser()
	if err != nil {
		logger.LogError(consts.ParserError, err)
		return err
	}

	return nil
}
