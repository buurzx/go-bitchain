package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// 24 is an arbitrary number, our goal is to have a target that
// takes less than 256 bits in memory.
// And we want the difference to be significant enough,
// but not too big, because the bigger the difference
// the more difficult it’s to find a proper hash.

// set it to 16 for speed
const targetBits = 16

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork function, we initialize a big.Int with the value of 1
// and shift it left by 256 - targetBits bits. 256 is the length of a
// SHA-256 hash in bits, and it’s SHA-256 hashing algorithm that we’re
// going to use. The hexadecimal representation of target is:
// 0000010000000000000000000000000000000000000000000000000000000000
// ================================================================
// The first hash (calculated on “I like donuts”) is bigger than the target,
// thus it’s not a valid proof of work. The second hash
// (calculated on “I like donutsca07ca”) is smaller than the target,
// thus it’s a valid proof.
// ================================================================
// You can think of a target as the upper boundary of a
// range: if a number (a hash) is lower than the boundary,
// it’s valid, and vice versa. Lowering the boundary will result
//  in fewer valid numbers, and thus, more difficult work
//  required to find a valid one.
// 0fac49161af82ed938add1d8725835cc123a1a87b1b196488360e58d4bfb51e3 "I like donuts"
// 0000010000000000000000000000000000000000000000000000000000000000 target
// 0000008b0f41ec78bab747864db66bcb9fb89920ee75f43fdaaeb5544f7f76ca "I like donutsca07ca"
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	// the counter
	nonce := 0

	fmt.Printf("Mining a new block")

	// “infinite” loop limited by maxNonce
	// this is done to avoid a possible overflow of nonce
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
