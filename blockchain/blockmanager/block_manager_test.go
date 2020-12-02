package blockmanager

import (
	"fmt"
	"github.com/reed/blockchain/store"
	"github.com/reed/crypto"
	"github.com/reed/database/leveldb"
	"github.com/reed/types"
	dbm "github.com/tendermint/tmlibs/db"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestBlockManager_calcFork(t *testing.T) {
	mainChain, forkChain, highBlock, curBlock := getMockBlock()

	fmt.Printf("highBlock %d %x \n", highBlock.Height, highBlock.GetHash())
	fmt.Printf("curBlock %d %x \n", curBlock.Height, curBlock.GetHash())

	bi := &BlockIndex{
		index: forkChain,
		main:  mainChain,
	}
	bm := &BlockManager{
		blockIndex: bi,
	}

	for _, block := range mainChain {
		if block == nil {
			continue
		}
		fmt.Printf("height:%d hash %x prev %x\n", block.Height, block.GetHash(), block.PrevBlockHash)
	}
	fmt.Println("====================================================================================")
	for _, block := range forkChain {
		fmt.Printf("height:%d hash %x prev %x\n", block.Height, block.GetHash(), block.PrevBlockHash)
	}
	fmt.Println("====================================================================================")

	reserves, discards, err := bm.calcFork(curBlock, highBlock)
	if err != nil {
		t.Error(err)
	}
	for _, block := range reserves {
		fmt.Printf("height:%d hash %x prev %x\n", block.Height, block.GetHash(), block.PrevBlockHash)
	}
	fmt.Println("====================================================================================")
	for _, block := range discards {
		fmt.Printf("height:%d hash %x prev %x\n", block.Height, block.GetHash(), block.PrevBlockHash)
	}
	fmt.Println("====================================================================================")

	if reserves[len(reserves)-1].PrevBlockHash != discards[len(discards)-1].PrevBlockHash {
		t.Error("calcFork error")
	}

}

func getMockBlock() ([]*types.Block, map[types.Hash]*types.Block, *types.Block, *types.Block) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	rn := r.Intn(9999999999)

	mainChain := make([]*types.Block, 20, 20)
	//for i := 0; i < 20; i++ {
	//	fmt.Println(i)
	//	mainChain[i]=nil
	//}
	forkChain := map[types.Hash]*types.Block{}

	prev := types.GenesisBlockHash()

	// common
	for i := 1; i < 10; i++ {
		flag := strconv.Itoa(i)
		mr := types.BytesToHash(crypto.Sha256([]byte("MerkleRootHash" + flag)))
		header := &types.BlockHeader{Height: uint64(i),
			PrevBlockHash:  prev,
			MerkleRootHash: mr,
			Timestamp:      uint64(time.Now().Unix()),
			Nonce:          uint64(rn),
			BigNumber:      *big.NewInt(int64(rn)),
			Version:        1,
		}

		block := &types.Block{BlockHeader: *header}
		mainChain[i] = block
		forkChain[block.GetHash()] = block
		prev = block.GetHash()
	}

	forkPre := prev

	var highBlock *types.Block
	for i := 10; i < 13; i++ {
		flag := strconv.Itoa(i)
		flag += "main-"
		mr := types.BytesToHash(crypto.Sha256([]byte("MerkleRootHash" + flag)))
		header := &types.BlockHeader{Height: uint64(i),
			PrevBlockHash:  prev,
			MerkleRootHash: mr,
			Timestamp:      uint64(time.Now().Unix()),
			Nonce:          uint64(rn),
			BigNumber:      *big.NewInt(int64(rn)),
			Version:        1,
		}

		block := &types.Block{BlockHeader: *header}
		mainChain[i] = block
		highBlock = block
		prev = block.GetHash()
	}

	var curBlock *types.Block
	for i := 10; i < 14; i++ {
		flag := strconv.Itoa(i)
		flag += "fork-"
		mr := types.BytesToHash(crypto.Sha256([]byte("MerkleRootHash" + flag)))
		header := &types.BlockHeader{Height: uint64(i),
			PrevBlockHash:  forkPre,
			MerkleRootHash: mr,
			Timestamp:      uint64(time.Now().Unix()),
			Nonce:          uint64(rn),
			BigNumber:      *big.NewInt(int64(rn)),
			Version:        1,
		}

		block := &types.Block{BlockHeader: *header}
		forkChain[block.GetHash()] = block
		forkPre = block.GetHash()
		curBlock = block
	}

	return mainChain, forkChain, highBlock, curBlock

}

func getStore() store.Store {
	return leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, os.Getenv("GOPATH")+"/src/github.com/reed/database/file/"))
}
