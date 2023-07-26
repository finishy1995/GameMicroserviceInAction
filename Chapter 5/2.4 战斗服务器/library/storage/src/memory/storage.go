package memory

import (
	"ProjectX/library/routine"
	"ProjectX/library/storage/core"
	"ProjectX/library/storage/src/tools"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type keyType uint8

const (
	Hash keyType = iota
	HashRange
)

const (
	DefaultTickTime = time.Second
	MinTickTime     = time.Millisecond * 20
)

type Storage struct {
	sync.RWMutex
	runFlag   bool
	maxLength int
	tick      time.Duration
	db        map[string]*table
}

type table struct {
	key             keyType
	preRefreshMutex sync.Mutex
	itemsMutex      sync.RWMutex
	items           map[string]*node
	preRefresh      map[string]bool // 预刷新队列，周期性更新
	head            node
	tail            node
}

type node struct {
	sync.Mutex
	front *node
	next  *node
	key   string
	value interface{}
}

func NewStorage(maxLength int, tick time.Duration) *Storage {
	stg := &Storage{
		runFlag:   true,
		maxLength: maxLength,
		db:        make(map[string]*table, 0),
	}
	if tick < MinTickTime {
		tick = DefaultTickTime
	}
	stg.tick = tick

	return stg
}

func (s *Storage) CreateTable(value interface{}) error {
	tableName := tools.GetStructOnlyName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, false)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}

	s.createTable(tableName, rangeKey)
	return nil
}

func (s *Storage) createTable(name string, rangeKey string) *table {
	s.RLock()
	if tb, ok := s.db[name]; ok {
		s.RUnlock()
		return tb
	}
	s.RUnlock()
	key := Hash
	if rangeKey != "" {
		key = HashRange
	}
	s.Lock()
	defer s.Unlock()
	if tb, ok := s.db[name]; ok {
		return tb
	}

	tb := &table{
		key:        key,
		items:      make(map[string]*node, 0),
		preRefresh: make(map[string]bool, 0),
		head:       node{},
		tail:       node{},
	}
	s.db[name] = tb
	s.db[name].head.next = &s.db[name].tail
	s.db[name].tail.front = &s.db[name].head

	// 最大长度有意义才执行
	if s.maxLength > 0 {
		err := routine.Run(false, func() {
			for {
				if s.runFlag {
					<-time.After(s.tick)
					s.process(tb)
				} else {
					return
				}
			}
		})
		if err != nil {
			panic(err)
		}
	}

	return s.db[name]
}

func (s *Storage) process(tb *table) {
	// 获取需要预刷新的队列
	tb.preRefreshMutex.Lock()
	length := len(tb.preRefresh)
	preRefreshSlice := make([]string, 0, length)
	for k := range tb.preRefresh {
		preRefreshSlice = append(preRefreshSlice, k)
	}
	tb.preRefresh = make(map[string]bool)
	tb.preRefreshMutex.Unlock()

	// 更新内存存储双向链表
	updateHead := &tb.head
	tb.itemsMutex.Lock()
	for _, k := range preRefreshSlice {
		if n, ok := tb.items[k]; ok {
			n.front.next = n.next
			n.next.front = n.front
			n.front = updateHead
			n.next = updateHead.next
			updateHead.next.front = n
			updateHead.next = n
			updateHead = n
		}
	}

	// 查询是否超载，超载则移除超载 node
	length = len(tb.items) - s.maxLength
	if length > 0 {
		updateTail := tb.tail.front
		for i := 0; i < length; i++ {
			delete(tb.items, updateTail.key)
			updateTail = updateTail.front
		}
		updateTail.next = &tb.tail
		tb.tail.front = updateTail
	}

	tb.itemsMutex.Unlock()
}

func (s *Storage) Create(value interface{}) error {
	tableName := tools.GetStructOnlyName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, false)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}

	tb := s.createTable(tableName, rangeKey)
	key := getRealKey(hashKey, rangeKey, tb.key, value)
	if key == "" {
		return core.ErrUnsupportedValueType
	}

	tb.itemsMutex.Lock()
	defer tb.itemsMutex.Unlock()
	// 检查对象是否存在
	if _, ok := tb.items[key]; ok {
		return core.ErrDuplicateKey
	}
	newNode := &node{
		front: &tb.head,
		next:  tb.head.next,
		value: value,
		key:   key,
	}
	tb.head.next.front = newNode
	tb.head.next = newNode
	tb.items[key] = newNode

	return nil
}

func getRealKey(hashKey string, rangeKey string, typ keyType, value interface{}) string {
	if typ == HashRange {
		if rangeKey == "" {
			return ""
		}
		return fmt.Sprintf("%v-%v",
			tools.GetFieldValueByName(value, hashKey),
			tools.GetFieldValueByName(value, rangeKey))
	}
	return fmt.Sprintf("%v", tools.GetFieldValueByName(value, hashKey))
}

func getRealKeyByValue(typ keyType, value ...interface{}) string {
	if typ == HashRange {
		return fmt.Sprintf("%v-%v", value[0], value[1])
	}
	return fmt.Sprintf("%v", value[0])
}

func (s *Storage) Delete(value interface{}, hash interface{}, args ...interface{}) error {
	tableName := tools.GetStructOnlyName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, false)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	var rangeValue interface{}
	if rangeKey != "" {
		if len(args) == 0 {
			return core.ErrMissingRangeValue
		}
		rangeValue = args[0]
	}

	tb := s.createTable(tableName, rangeKey)
	key := getRealKeyByValue(tb.key, hash, rangeValue)
	if key == "" {
		return core.ErrUnsupportedValueType
	}

	tb.itemsMutex.Lock()
	defer tb.itemsMutex.Unlock()
	if delNode, ok := tb.items[key]; ok {
		delNode.front.next = delNode.next
		delNode.next.front = delNode.front
		delete(tb.items, key)
	}
	return nil
}

func (s *Storage) Save(value interface{}) error {
	tableName := tools.GetStructName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, false)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	_, err := tools.TrySetStructVersion(value)
	if err != nil {
		return err
	}

	tb := s.createTable(tableName, rangeKey)
	key := getRealKey(hashKey, rangeKey, tb.key, value)
	if key == "" {
		return core.ErrUnsupportedValueType
	}

	tb.itemsMutex.RLock()
	getNode, ok := tb.items[key]
	tb.itemsMutex.RUnlock()
	tb.preRefreshMutex.Lock()
	tb.preRefresh[key] = true
	tb.preRefreshMutex.Unlock()

	if ok {
		original := reflect.ValueOf(value).Elem()
		cpy := reflect.New(original.Type())
		err = tools.DeepCopy(value, cpy.Interface())
		getNode.value = cpy.Elem().Interface()
	}

	return err
}

func (s *Storage) First(value interface{}, hash interface{}, args ...interface{}) error {
	tableName := tools.GetStructName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, false)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	var rangeValue interface{}
	if rangeKey != "" {
		if len(args) == 0 {
			return core.ErrMissingRangeValue
		}
		rangeValue = args[0]
	}

	tb := s.createTable(tableName, rangeKey)
	key := getRealKeyByValue(tb.key, hash, rangeValue)
	if key == "" {
		return core.ErrUnsupportedValueType
	}

	tb.itemsMutex.RLock()
	getNode, ok := tb.items[key]
	tb.itemsMutex.RUnlock()

	if ok {
		tb.preRefreshMutex.Lock()
		tb.preRefresh[key] = true
		tb.preRefreshMutex.Unlock()
		err := tools.DeepCopy(getNode.value, value)
		if err != nil {
			return err
		}
	} else {
		return core.ErrNotFound
	}
	return nil
}

// Find 支持 > ; < ; >= ; <= ; <> ; = ; and ; or ; () ; not
// 本地存储需要新增，可以按照上述计算符添加
func (s *Storage) Find(value interface{}, limit int64, expr string, args ...interface{}) error {
	tableName := tools.GetSliceStructName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, false)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}

	tb := s.createTable(tableName, rangeKey)
	if limit == 0 {
		limit--
	}
	var count int64 = 0
	var nod *exprNode
	if expr != "" {
		nod = getExprRoot(expr)
		if nod == nil {
			return core.ErrUnsupportedExprType
		}
	}
	var slc []interface{}
	tb.itemsMutex.RLock()
	defer tb.itemsMutex.RUnlock()

	for _, item := range tb.items {
		if nod.calculate(item.value, args) {
			slc = append(slc, item.value)

			count++
			if limit == count {
				break
			}
		}
	}

	return tools.DeepCopy(slc, value)
}
