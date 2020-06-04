package dav

import (
	"github.com/mknentwich/core/context"
	"golang.org/x/net/webdav"
	"strings"
)

type BillCollectionNode struct {
	*BasicNode
	Payed bool
}

func (b *BillCollectionNode) Name() string {
	name := context.Conf.DavPaths.UnpayedBills
	if b.Payed {
		name = context.Conf.DavPaths.PayedBills
	}
	splitName := strings.Split(name, separator)
	return splitName[len(splitName)-1]
}

func newBillCollectionNode(payed bool) *BillCollectionNode {
	node := &BillCollectionNode{Payed: payed}
	node.BasicNode = newBasicNode(node.Name())
	return node
}

type BillNode struct {
	*BasicNode
}

func (b *BillCollectionNode) File() webdav.File {
	return &BasicFile{node: b.BasicNode}
}

type BillFile struct {
}

type BillCollectionFile struct {
}

type billFilter struct {
	payed *bool
	month *int
	year  *int
	id    *int
}
