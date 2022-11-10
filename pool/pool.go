package pool

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/trees/redblacktree"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	count     int64
	path      string
	lock      sync.Mutex
	tokenTree *redblacktree.Tree
}

func keyComparator(left, right interface{}) int {

	lefts := strings.Split(left.(string), "-")
	rights := strings.Split(right.(string), "-")

	leftValue, _ := strconv.ParseInt(lefts[0], 10, 64)
	rightValue, _ := strconv.ParseInt(rights[0], 10, 64)

	res := leftValue - rightValue

	if res != 0 {
		return int(res)
	} else {
		leftValue, _ = strconv.ParseInt(lefts[1], 10, 64)
		rightValue, _ = strconv.ParseInt(rights[1], 10, 64)

		return int(leftValue - rightValue)
	}

}

func (pool *Pool) getWalk() filepath.WalkFunc {

	var token Token
	var data []byte

	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			SugarLogger.Error(err)
			return err
		}

		if info.IsDir() {
			return err
		}

		data, err = os.ReadFile(path)
		if err != nil {
			SugarLogger.Error(err)
			return err
		}

		err = json.Unmarshal(data, &token)
		if err != nil {
			SugarLogger.Error(err)
			return err
		}

		pool.Offer(token)
		return nil
	}
}

func (pool *Pool) getRefreshWalk() filepath.WalkFunc {

	var token Token
	var data []byte
	now := time.Now().UnixMilli()

	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			SugarLogger.Error(err)
			return err
		}

		if info.IsDir() {
			return err
		}

		data, err = os.ReadFile(path)
		if err != nil {
			SugarLogger.Error(err)
			return err
		}

		err = json.Unmarshal(data, &token)
		if err != nil {
			SugarLogger.Error(err)
			return err
		}

		if token.ExpireTimestamp <= now {
			os.Remove(path)
		}
		return nil
	}
}

func (pool *Pool) runTask() {

	ticker := time.NewTicker(10 * time.Second)

	go func(pool *Pool, t *time.Ticker) {
		for {
			<-t.C
			SugarLogger.Info("here")
			pool.refreshPool()
		}

		SugarLogger.Info("error here")
	}(pool, ticker)

}

func New(path string) *Pool {

	tree := redblacktree.NewWith(keyComparator)

	res := &Pool{
		count:     0,
		path:      path,
		lock:      sync.Mutex{},
		tokenTree: tree,
	}

	res.runTask()
	filepath.Walk(path, res.getWalk())
	return res

}

func (pool *Pool) refreshPool() {

	SugarLogger.Info("refresh here !")
	err := filepath.Walk(pool.path, pool.getRefreshWalk())
	if err != nil {
		SugarLogger.Error(err)
	}

	key := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), 0)
	defer pool.lock.Unlock()
	pool.lock.Lock()
	for {
		token, ok := pool.tokenTree.Floor(key)
		if !ok {
			break
		} else {
			pool.tokenTree.Remove(token.Key)
		}
	}

}

func (pool *Pool) addCount() int64 {
	atomic.AddInt64(&pool.count, 1)
	return pool.count
}

func (pool *Pool) Offer(token Token) {
	var keyTimestamp int64
	var current int64

	if token.ExpireTimestamp == 0 {
		now := time.Now()
		keyTimestamp = now.UnixMilli() + 540*1000
		current = pool.addCount()
		token.ExpireTimestamp = keyTimestamp
		token.Id = current

		fileName := fmt.Sprintf("%s\\%s-%d", pool.path, now.Format("15_04_05"), current)
		data, err := token.toJson()
		if err == nil {
			os.WriteFile(fileName, data, 0666)
		} else {
			SugarLogger.Error(err)
		}
	} else {
		keyTimestamp = token.ExpireTimestamp
		current = token.Id
	}

	defer pool.lock.Unlock()
	pool.lock.Lock()
	pool.tokenTree.Put(fmt.Sprintf("%d-%d", keyTimestamp, current), token)
}

func (pool *Pool) Get(timestamp int64) *Token {
	key := fmt.Sprintf("%d-%d", timestamp, 0)
	defer pool.lock.Unlock()
	pool.lock.Lock()
	res, ok := pool.tokenTree.Ceiling(key)
	if ok {
		pool.tokenTree.Remove(res.Key)
		return res.Value.(*Token)
	}

	return nil
}
