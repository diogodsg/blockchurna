package blockchain

type Blockchain struct {
	Blocks []*Block
}

var BC *Blockchain = NewBlockchain()

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, "", "0")
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

func (bc *Blockchain) CreateBlock(payload string) *Block {
	previousBlock := bc.Blocks[len(bc.Blocks) - 1]
	newBlock := NewBlock(previousBlock.Index+1, payload, previousBlock.Id)
	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock
}

func (bc *Blockchain) AddBlock(block *Block) *Block {
	latestBlock := bc.GetLatestBlock()
	if block.Index <= latestBlock.Index {
		return block
	}
	bc.Blocks = append(bc.Blocks, block)

	return block
}

func (bc *Blockchain) ReplaceBlockchain(blocks []*Block) []*Block {
	bc.Blocks = blocks
	return blocks
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks) - 1]
}


func (b *Blockchain) GetBlocksAfterIndex(index int) []*Block {
	return b.Blocks
	// if index >= len(b.Blocks)-1 {
    //     return []*Block{}
    // }
    // return b.Blocks[index+1:]
}