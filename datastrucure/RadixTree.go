package datastrucure

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
		// If prefix is higher than supported compression we need to split
		if prefixLen > 64 {
			d.root = &RadixTreeNode{compressed: high, compressionLength: 64}
			if (low >> 63) == 1 {
				d.root.right = &RadixTreeNode{compressed: low << 1, compressionLength: prefixLen - 64, popId: popNumber}
			} else {
				d.root.left = &RadixTreeNode{compressed: low << 1, compressionLength: prefixLen - 64, popId: popNumber}
			}
		} else {
			d.root = &RadixTreeNode{compressed: high, compressionLength: prefixLen, popId: popNumber}
		}
	} else {
		d.root = insertNonRootNode(d.root, high, low, prefixLen, popNumber)
	}
}

func insertNonRootNode(node *RadixTreeNode, high uint64, low uint64, prefixLen uint8, popNumber uint8) *RadixTreeNode
