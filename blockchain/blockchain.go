package blockchain

import "fmt"

type Blockchain struct {
	Blocks []*Block
}

var BC *Blockchain = NewBlockchain()

func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{}}
}

func (bc *Blockchain) CreateBlock(payload Payload) *Block {
	previousNode := ""
	index := 0

	if len(bc.Blocks) > 0 {
		previousBlock :=  bc.Blocks[len(bc.Blocks) - 1]
		previousNode = previousBlock.Id
		index = previousBlock.Index + 1 
	} 
	err := ValidatePayload(payload)

	if err != nil {
		fmt.Println("Invalid Block")
		return nil
	}

	newBlock := NewBlock(index, payload, previousNode)
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