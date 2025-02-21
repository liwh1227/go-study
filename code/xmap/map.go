package xmap

import "fmt"

// 键值对节点
type entry struct {
	key   string
	value interface{}
	next  *entry
}

type HashTable struct {
	buckets []*entry // 桶数组，每个元素是一个链表的头指针
	size    int      // 哈希表大小（桶的数量）
}

// 创建哈希表
func NewHashTable(size int) *HashTable {
	return &HashTable{
		buckets: make([]*entry, size),
		size:    size,
	}
}

// 哈希函数
func (ht *HashTable) hash(key string) int {
	hash := 0

	for _, char := range key {
		hash = (hash*31 + int(char)) % ht.size
	}

	return hash
}

func (ht *HashTable) Insert(key string, value interface{}) {
	// 计算下标
	index := ht.hash(key)

	e := &entry{
		key:   key,
		value: value,
		next:  nil,
	}

	if ht.buckets[index] == nil {
		// 说明该桶无冲突
		ht.buckets[index] = e
	} else {
		current := ht.buckets[index]
		for current != nil {
			if current.key == key {
				// key存在，则更新
				current.value = value
				return
			}
			if current.next == nil {
				// 说明链表尾，需要添加
				break
			}
			current = current.next
		}
		// 键不存在，插入新元素
		current.next = e
	}
}

// 获取
// 1. 该桶不存在该key的hash
// 2. 存在该key的hash，但是没有该key生成hash
func (ht *HashTable) Get(key string) (interface{}, bool) {
	index := ht.hash(key)
	if ht.buckets[index] == nil {
		// 没有此key对应的值
		return nil, false
	}

	current := ht.buckets[index]
	for current != nil {
		if current.key == key {
			return current.value, true
		}
		current = current.next
	}

	return nil, false
}

// 删除
func (ht *HashTable) Del(key string) {
	index := ht.hash(key)
	if ht.buckets[index] == nil {
		return
	}

	if ht.buckets[index].key == key {
		ht.buckets[index] = ht.buckets[index].next
		return
	}

	prev := ht.buckets[index]
	current := prev.next

	for current != nil {
		if current.key == key {
			prev.next = current.next
			break
		}
		prev = current
		current = current.next
	}
}

func TestHashTable() {
	ht := NewHashTable(5)

	ht.Insert("h", 10)
	ht.Insert("e", 100)

	v, ok := ht.Get("h")
	if ok {
		fmt.Println(v)
	}

	v, ok = ht.Get("e")
	if ok {
		fmt.Println(v)
	}

	v, ok = ht.Get("l")
	if ok {
		fmt.Println(v)
	}

	ht.Del("h")
	v, ok = ht.Get("h")
	if ok {
		fmt.Println(v)
	}
}

/*
A map is just a hash table.
The data is arranged into an array of buckets.
Each bucket contains up to 8 key/elem pairs.
The low-order bits of the hash are used to select a bucket.
Each bucket contains a few high-order bits of each hash to distinguish the entries within a single bucket. */

//Map iterators walk through the array of buckets and
//return the keys in walk order (bucket #, then overflow
//chain order, then bucket index).  To maintain iteration
//semantics, we never move keys within their bucket (if
//we did, keys might be returned 0 or 2 times).  When
//growing the table, iterators remain iterating through the
//old table and must check the new table if the bucket
//they are iterating through has been moved ("evacuated")
//to the new table.
