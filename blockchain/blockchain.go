package blockchain

type Blockchain struct {
	Blocks []*Block
}

var BC *Blockchain = NewBlockchain()

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, "", "0")
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

func (bc *Blockchain) AddBlock(payload string) *Block {
	previousBlock := bc.Blocks[len(bc.Blocks) - 1]
	newBlock := NewBlock(previousBlock.Index+1, payload, previousBlock.Id)
	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks) - 1]
}

