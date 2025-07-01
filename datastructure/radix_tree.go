package datastructure

import "github.com/vbaz365/cdn/utils"

type RadixTreeNode struct {
    compression        uint64
    compressionLen uint8
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
    }

    insertNonRootNode(d.root, high, low, prefixLen, popNumber)
}

func insertNonRootNode(node *RadixTreeNode, high uint64, low uint64, prefixLen uint8, popNumber uint8){
	// Go through nodes until we add everything
    var index uint8 = 0
    for index < prefixLen {
		var firstBit uint8 = utils.GetFirstBit(utils.GetBitsAfterIndex(high, low, index, prefixLen))
		remainingBitsCount := prefixLen - index
		// Check the right/left child of current node
		if firstBit == 1 && node.right == nil{
            node.right = &RadixTreeNode{}
			newNodeLogic(node.right, &index, remainingBitsCount, high, low, prefixLen, popNumber)
            node = node.right
            continue
		} else if firstBit == 0 && node.left == nil{
            node.left = &RadixTreeNode{}
			newNodeLogic(node.left, &index, remainingBitsCount, high, low, prefixLen, popNumber)
            node = node.left
            continue
		}
        // If there is a child
        if firstBit == 1{
            node = node.right
        } else{
            node = node.left
        }

        var bitMatchLength uint8 = utils.CommonPrefixLen(utils.GetBitsAfterIndex(high, low, index, prefixLen), node.compression)
        if bitMatchLength < node.compressionLen{
            splitNodeLogic(node, bitMatchLength)
        }
        index += bitMatchLength
    }
}

func newNodeLogic(toBeAdded *RadixTreeNode, index *uint8, remainingBitsCount uint8, high uint64, low uint64, prefixLen uint8, popNumber uint8){
    // If the remaining unsaved bits can be compressed into one node
    if remainingBitsCount <= 64 && *index < 64 && prefixLen > 64{
        compressed, compressedLen := utils.CombineSmallCompressions(high, 64 - *index, low, prefixLen - *index)
        updateNode(toBeAdded, compressed, compressedLen, popNumber)
        *index += toBeAdded.compressionLen
    } else if prefixLen > 64 && *index < 64{
        updateNode(toBeAdded, utils.GetBitsAfterIndex(high, low, *index, prefixLen), 64 - *index, 0)
        *index += toBeAdded.compressionLen
    } else if prefixLen > 64{
        updateNode(toBeAdded, utils.GetBitsAfterIndex(high, low, *index, prefixLen), prefixLen - *index, popNumber)
        *index += toBeAdded.compressionLen
    }
}

func splitNodeLogic(toBeSplitted *RadixTreeNode, splitIndex uint8){
    highCompression, lowCompression := utils.SplitCompression(toBeSplitted.compression, splitIndex)
    var highCompressionLen uint8 = splitIndex
    var lowCompressionLen uint8 = toBeSplitted.compressionLen - splitIndex
    tempRight := toBeSplitted.right
    tempLeft := toBeSplitted.left
    if utils.GetFirstBit(lowCompression) == 1{
        toBeSplitted.right = &RadixTreeNode{compression: lowCompression, compressionLen: lowCompressionLen, popId: toBeSplitted.popId, right: tempRight, left: tempLeft}
        toBeSplitted.left = nil
    } else {
        toBeSplitted.left = &RadixTreeNode{compression: lowCompression, compressionLen: lowCompressionLen, popId: toBeSplitted.popId, right: tempRight, left: tempLeft}
        toBeSplitted.right = nil
    }
    toBeSplitted.compression = highCompression
    toBeSplitted.compressionLen = highCompressionLen
    toBeSplitted.popId = 0
}

func updateNode(node *RadixTreeNode, compressed uint64, compressedLen uint8, popNumber uint8){
    node.compression = compressed
    node.compressionLen = compressedLen
    node.popId = popNumber
}