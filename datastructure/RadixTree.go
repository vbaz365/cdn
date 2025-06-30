package datastructure

import "github.com/vbaz365/cdn/utils"

type RadixTreeNode struct {
    compressed        uint64
    compressionLength uint8
    popId             uint8
    right             *RadixTreeNode
    left              *RadixTreeNode
}

type Data struct {
    root *RadixTreeNode
}

func (d *Data) insertRadixNode(high uint64, low uint64, prefixLen uint8, popNumber uint8) {
    // If first node
    if d.root == nil {
        d.root = &RadixTreeNode{}
    } else {
        d.root = insertNonRootNode(d.root, high, low, prefixLen, popNumber, 0)
    }
}

func insertNonRootNode(node *RadixTreeNode, high uint64, low uint64, prefixLen uint8, popNumber uint8, index uint8) *RadixTreeNode {
	// Go through nodes until we add everything
    for index < prefixLen {
		var firstBit uint8 = utils.GetFirstBit(utils.GetBitsAfterIndex(high, low, index))
		remainingBitsCount := prefixLen - index
		// Check the right/left child of current node
		if(firstBit == 1 && node.right == nil){
			newNodeLogic(node, node.right, &index)
		} else if(node.left == nil){
			newNodeLogic(node, node.left, &index)
		}
    }
}

func newNodeLogic(current *RadixTreeNode, toBeAdded *RadixTreeNode, index *uint8){
    // If the remaining unsaved bits can be compressed into one node
    if(remainingBitsCount <= 64 && index < 64 && prefixLen > 64){
        compressed, compressedLen := utils.CombineSmallCompressions(high, 64 - index, low, prefixLen - index)
        toBeAdded = &RadixTreeNode{compression: compressed, compressionLen: compressedLen, popId: popNumber}
        current = toBeAdded
        index += current.compressionLen
    } else if (prefixLen > 64 && index < 64){
        toBeAdded = &RadixTreeNode{compression: utils.GetBitsAfterIndex(high, low, index), compressionLen: 64 - index}
        current = toBeAdded
        index += current.compressionLen
    } else if (prefixLen > 64){
        toBeAdded = &RadixTreeNode{compression: utils.GetBitsAfterIndex(high, low, index), compressionLen: prefixLen - index}
        current = toBeAdded
        index += current.compressionLen
    }
}