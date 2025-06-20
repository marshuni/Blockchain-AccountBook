package db

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/pow"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/tx"
)

func init() {
	gob.Register(&pow.Block{})
	gob.Register(&tx.Transaction{})
}

// 区块存储桶名
var blocksBucket = []byte("blocks")
var lastHashKey = []byte("lastHash")

type DB struct {
	db *bolt.DB
}

// 打开数据库
func OpenDB(path string) (*DB, error) {
	database, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// 初始化存储桶
	err = database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(blocksBucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &DB{db: database}, nil
}

// 存储区块
func (d *DB) PutBlock(hash []byte, block interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(block); err != nil {
		return err
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		return b.Put(hash, buf.Bytes())
	})
}

// 读取区块
func (d *DB) GetBlock(hash []byte, block interface{}) error {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		data := b.Get(hash)
		if data == nil {
			return errors.New("block not found")
		}
		dec := gob.NewDecoder(bytes.NewReader(data))
		return dec.Decode(block)
	})
}

// 获取最后一个区块的哈希值
func (d *DB) GetLastHash() ([]byte, error) {
	var lastHash []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		lastHash = b.Get(lastHashKey)
		if lastHash == nil {
			return errors.New("last hash not found")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return lastHash, nil
}

func (d *DB) UpdateLastHash(hash []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		return b.Put(lastHashKey, hash)
	})
}

// 关闭数据库
func (d *DB) Close() error {
	return d.db.Close()
}
