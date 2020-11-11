package blockchain

type Chain struct {
	Store  Store
	Txpool *Txpool
}
