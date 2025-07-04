package datastructure

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/vbaz365/cdn/utils"
)

const InvalidPopID uint16 = 65535

type RadixTreeNode struct {
	compression    uint64
	compressionLen uint8
	popId          uint16
	right          *RadixTreeNode
	left           *RadixTreeNode
}

type Data struct {
	root *RadixTreeNode
}

// InsertRadixNode inserts a new node into data
func (d *Data) InsertRadixNode(high uint64, low uint64, prefixLen uint8, popNumber uint16) {
	// If first node
	if d.root == nil {
		d.root = &RadixTreeNode{}
	}

	insertNonRootNode(d.root, high, low, prefixLen, popNumber)
}

// insertNonRootNode is a helper function to insert a new node into the data
func insertNonRootNode(node *RadixTreeNode, high uint64, low uint64, prefixLen uint8, popNumber uint16) {
	// Go through nodes until we add everything
	var index uint8 = 0
	for index < prefixLen {
		var firstBit uint8 = utils.GetFirstBit(utils.GetBitsAfterIndex(high, low, index))
		remainingBitsCount := prefixLen - index
		// Check the right/left child of current node
		if firstBit == 1 && node.right == nil {
			node.right = &RadixTreeNode{}
			newNodeLogic(node.right, &index, remainingBitsCount, high, low, prefixLen, popNumber)
			node = node.right
			continue
		} else if firstBit == 0 && node.left == nil {
			node.left = &RadixTreeNode{}
			newNodeLogic(node.left, &index, remainingBitsCount, high, low, prefixLen, popNumber)
			node = node.left
			continue
		}
		// If there is a child
		if firstBit == 1 {
			node = node.right
		} else {
			node = node.left
		}

		// Check how many bits match with the child
		var bitMatchLength uint8 = utils.CommonPrefixLen(utils.GetBitsAfterIndex(high, low, index), node.compression)
		if bitMatchLength > remainingBitsCount {
			bitMatchLength = remainingBitsCount
		}

		// If less than required bits match split the node otherwise continue down that node
		if bitMatchLength < node.compressionLen {
			splitNodeLogic(node, bitMatchLength)
			index += bitMatchLength
			if index == prefixLen {
				node.popId = popNumber
			}
		} else {
			index += node.compressionLen
		}

	}
}

// newNodeLogic creates a new node with as many bits of IPv6 as possible
func newNodeLogic(toBeAdded *RadixTreeNode, index *uint8, remainingBitsCount uint8, high uint64, low uint64, prefixLen uint8, popNumber uint16) {
	// If the remaining unsaved bits can be compressed into one node
	if remainingBitsCount <= 64 && *index < 64 && prefixLen > 64 {
		compressed, compressedLen := utils.CombineSmallCompressions(high, 64-*index, low, prefixLen-*index)
		updateNode(toBeAdded, compressed, compressedLen, popNumber)
		*index += toBeAdded.compressionLen
	} else if prefixLen > 64 && *index < 64 {
		updateNode(toBeAdded, utils.GetBitsAfterIndex(high, low, *index), 64-*index, InvalidPopID)
		*index += toBeAdded.compressionLen
	} else {
		updateNode(toBeAdded, utils.GetBitsAfterIndex(high, low, *index), prefixLen-*index, popNumber)
		*index += toBeAdded.compressionLen
	}
}

// splitNodeLogic splits a single node into two connected nodes
func splitNodeLogic(toBeSplitted *RadixTreeNode, splitIndex uint8) {
	highCompression, lowCompression := utils.SplitCompression(toBeSplitted.compression, splitIndex)
	var highCompressionLen uint8 = splitIndex
	var lowCompressionLen uint8 = toBeSplitted.compressionLen - splitIndex
	tempRight := toBeSplitted.right
	tempLeft := toBeSplitted.left

	// Set lower compression, children and popId into a new child
	if utils.GetFirstBit(lowCompression) == 1 {
		toBeSplitted.right = &RadixTreeNode{compression: lowCompression, compressionLen: lowCompressionLen, popId: toBeSplitted.popId, right: tempRight, left: tempLeft}
		toBeSplitted.left = nil
	} else {
		toBeSplitted.left = &RadixTreeNode{compression: lowCompression, compressionLen: lowCompressionLen, popId: toBeSplitted.popId, right: tempRight, left: tempLeft}
		toBeSplitted.right = nil
	}

	// Set the upper compression to the node and reset the popId
	toBeSplitted.compression = highCompression
	toBeSplitted.compressionLen = highCompressionLen
	toBeSplitted.popId = InvalidPopID
}

// updateNode is a helper function to set node variables
func updateNode(node *RadixTreeNode, compressed uint64, compressedLen uint8, popNumber uint16) {
	node.compression = compressed
	node.compressionLen = compressedLen
	node.popId = popNumber
}

// Route returns pop ID and scope prefix-length
func (d *Data) Route(ecs *net.IPNet) (pop uint16, scope int) {
	if d.root == nil || ecs == nil {
		return 0, 0
	}

	// Check if we have valid IPv6
	ip := ecs.IP.To16()
	if ip == nil || len(ip) != net.IPv6len {
		return 0, 0
	}

	prefixLen, _ := ecs.Mask.Size()

	// Get higher and lower half of the IPv6 represented as 2 uint64s
	high := binary.BigEndian.Uint64(ip[:8])
	low := binary.BigEndian.Uint64(ip[8:])

	node := d.root
	var matchedPop uint16 = 0
	var matchedPrefix int = 0
	var index uint16 = 0

	firstBit := utils.GetFirstBit(high)
	if firstBit == 1 {
		node = node.right
	} else {
		node = node.left
	}

	// Traverse the nodes and save the most recent valid popId and scope
	for node != nil && index < uint16(prefixLen) {
		compressLen := node.compressionLen
		bitsAtIndex := utils.GetBitsAfterIndex(high, low, uint8(index))
		commonLen := utils.CommonPrefixLen(bitsAtIndex, node.compression)

		if commonLen < compressLen {
			break
		}

		index += uint16(compressLen)

		if node.popId != InvalidPopID {
			fmt.Println(node.popId)
			matchedPop = node.popId
			matchedPrefix = int(index)
		}

		nextBit := utils.GetFirstBit(utils.GetBitsAfterIndex(high, low, uint8(index)))
		if nextBit == 1 {
			node = node.right
		} else {
			node = node.left
		}
	}

	return matchedPop, matchedPrefix
}
