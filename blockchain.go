package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const blocksBucket = "blocks"
const dbFile = "blockchain_%s.db"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock to blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	// get the last block hash from the DB to use it to mine a new block hash.
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	// mining a new block
	newBlock := NewBlock(data, lastHash)
	// save its serialized representation into the DB
	// and update the l key, which now stores the new block’s hash.
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			log.Panic(err)
		}
		if err = b.Put([]byte("l"), newBlock.Hash); err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
}

// NewGenesisBlock must be the first block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain with Genesis block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}
