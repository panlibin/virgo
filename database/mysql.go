package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	logger "github.com/panlibin/vglog"
	"github.com/panlibin/virgo"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ErrCallbackIsNil 无回调
var ErrCallbackIsNil = errors.New("database query callback is nil")

const defaultQueryChannelSize = 1024
const (
	queryTypeQuery int32 = iota
	queryTypeQueryRow
	queryTypeExec
)

type mysqlQueryContext struct {
	query        string
	args         []interface{}
	queryType    int32
	callbackChan chan []interface{}
}

type mysqlInstance struct {
	db              *sql.DB
	queryChan       chan *mysqlQueryContext
	wg              sync.WaitGroup
	aliveTicker     *time.Ticker
	cancelAliveCtx  context.Context
	cancelAliveFunc context.CancelFunc
}

func (m *mysqlInstance) open(dsn string) (err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	m.db = db
	m.queryChan = make(chan *mysqlQueryContext, defaultQueryChannelSize)

	go m.run()
	go m.keepAlive()

	return
}

func (m *mysqlInstance) close() {
	if m.queryChan != nil {
		m.queryChan <- nil
	}

	if m.aliveTicker != nil {
		m.aliveTicker.Stop()
		m.aliveTicker = nil
	}
	if m.cancelAliveFunc != nil {
		m.cancelAliveFunc()
		m.cancelAliveFunc = nil
	}

	m.wg.Wait()
	if m.db != nil {
		m.db.Close()
	}
}

func (m *mysqlInstance) run() {
	m.wg.Add(1)
	defer m.wg.Done()
	for queryCtx := range m.queryChan {
		if queryCtx == nil {
			break
		}
		var ret interface{}
		var err error
		switch queryCtx.queryType {
		case queryTypeQuery:
			ret, err = m.db.Query(queryCtx.query, queryCtx.args...)
		case queryTypeQueryRow:
			ret = m.db.QueryRow(queryCtx.query, queryCtx.args...)
		case queryTypeExec:
			ret, err = m.db.Exec(queryCtx.query, queryCtx.args...)
		default:
			continue
		}

		if err != nil {
			logger.Errorf("%v", err)
			logger.Errorf(queryCtx.query+"; "+strings.Repeat("%v\t", len(queryCtx.args)), queryCtx.args...)
		}

		if queryCtx.callbackChan != nil {
			queryCtx.callbackChan <- []interface{}{ret, err}
		}
	}
	close(m.queryChan)
}

func (m *mysqlInstance) keepAlive() {
	m.cancelAliveCtx, m.cancelAliveFunc = context.WithCancel(context.Background())
	m.aliveTicker = time.NewTicker(time.Minute * 10)
	bQuit := false
	for !bQuit {
		select {
		case <-m.aliveTicker.C:
			m.db.Ping()
		case <-m.cancelAliveCtx.Done():
			bQuit = true
		}
	}
}

func (m *mysqlInstance) addQuery(queryCtx *mysqlQueryContext) {
	m.queryChan <- queryCtx
}

// Mysql 数据库管理对象
type Mysql struct {
	arrDb []*mysqlInstance
}

// NewMysql 新建
func NewMysql() *Mysql {
	pDb := new(Mysql)
	return pDb
}

// Open 连接数据库
func (m *Mysql) Open(dsn string, instNum int32) error {
	m.arrDb = make([]*mysqlInstance, instNum)
	var err error
	for i := int32(0); i < instNum; i++ {
		pDbInst := new(mysqlInstance)
		err = pDbInst.open(dsn)
		if err != nil {
			break
		}
		m.arrDb[i] = pDbInst
	}

	if err != nil {
		m.Close()
	}

	return err
}

// Close 关闭数据库连接
func (m *Mysql) Close() {
	if m.arrDb != nil {
		for _, pDb := range m.arrDb {
			if pDb != nil {
				pDb.close()
			}
		}
	}
}

// Query 查询多行
func (m *Mysql) Query(dbIdx uint32, query string, args ...interface{}) (rows *sql.Rows, err error) {
	callbackChan := m.pushOperator(dbIdx, query, args, queryTypeQuery, true)

	ret := <-callbackChan
	close(callbackChan)

	if ret[0] != nil {
		rows = ret[0].(*sql.Rows)
	}
	if ret[1] != nil {
		err = ret[1].(error)
	}

	return
}

// QueryRow 查询一行
func (m *Mysql) QueryRow(dbIdx uint32, query string, args ...interface{}) (row *sql.Row) {
	callbackChan := m.pushOperator(dbIdx, query, args, queryTypeQueryRow, true)

	ret := <-callbackChan
	close(callbackChan)

	return ret[0].(*sql.Row)
}

// Exec 执行
func (m *Mysql) Exec(dbIdx uint32, query string, args ...interface{}) (res sql.Result, err error) {
	callbackChan := m.pushOperator(dbIdx, query, args, queryTypeExec, true)

	ret := <-callbackChan
	close(callbackChan)

	if ret[0] != nil {
		res = ret[0].(sql.Result)
	}
	if ret[1] != nil {
		err = ret[1].(error)
	}

	return
}

// AsyncQuery 查询多行,回调
func (m *Mysql) AsyncQuery(p virgo.IProcedure, cb func([]interface{}), dbIdx uint32, query string, args ...interface{}) error {
	if p == nil || cb == nil {
		return ErrCallbackIsNil
	}

	callbackChan := m.pushOperator(dbIdx, query, args, queryTypeQuery, true)

	go func() {
		ret := <-callbackChan
		close(callbackChan)
		p.SyncTask(cb, ret...)
	}()

	return nil
}

// AsyncQueryRow 查询一行,回调
func (m *Mysql) AsyncQueryRow(p virgo.IProcedure, cb func([]interface{}), dbIdx uint32, query string, args ...interface{}) error {
	if p == nil || cb == nil {
		return ErrCallbackIsNil
	}

	callbackChan := m.pushOperator(dbIdx, query, args, queryTypeQueryRow, true)

	go func() {
		ret := <-callbackChan
		close(callbackChan)
		p.SyncTask(cb, ret...)
	}()

	return nil
}

// AsyncExec 执行,回调
func (m *Mysql) AsyncExec(p virgo.IProcedure, cb func([]interface{}), dbIdx uint32, query string, args ...interface{}) {
	needCb := (p != nil && cb != nil)

	callbackChan := m.pushOperator(dbIdx, query, args, queryTypeExec, needCb)

	if needCb {
		go func() {
			ret := <-callbackChan
			close(callbackChan)
			p.SyncTask(cb, ret...)
		}()
	}
}

func (m *Mysql) pushOperator(dbIdx uint32, query string, args []interface{}, queryType int32, needCb bool) chan []interface{} {
	dbCount := uint32(len(m.arrDb))
	if dbIdx >= dbCount {
		dbIdx %= dbCount
	}

	db := m.arrDb[dbIdx]
	queryCtx := new(mysqlQueryContext)
	queryCtx.query = query
	queryCtx.args = args
	queryCtx.queryType = queryType
	var callbackChan chan []interface{}
	if needCb {
		callbackChan = make(chan []interface{}, 1)
		queryCtx.callbackChan = callbackChan
	}

	db.addQuery(queryCtx)

	return callbackChan
}