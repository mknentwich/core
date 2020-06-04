package dav

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"golang.org/x/net/webdav"
	"math/rand"
	"os"
	"strings"
	"time"
)

var tree *BasicNode

type PhantomNode interface {
	Children() []PhantomNode
	File() webdav.File
	Name() string
	Subset(string) PhantomNode
}

type BasicNode struct {
	children []PhantomNode
	name     string
}

func (b *BasicNode) Children() []PhantomNode {
	//np := *b
	//n, err := np.(*CategoryNode)
	return b.children
}

func (b *BasicNode) File() webdav.File {
	return &BasicFile{
		node: b,
	}
}

func (b *BasicNode) Name() string {
	return b.name
}

func (b *BasicNode) Subset(path string) PhantomNode {
	if path == "" {
		return b
	}
	name := strings.Split(path, separator)[0]
	for _, child := range b.Children() {
		if child.Name() == name {
			if len(name) == len(path) {
				return child
			}
			return child.Subset(path[len(name)+1:])
		}
	}
	return nil
}

//appends the given node
//the resulting path of the node will be the given path plus the nodes name
//if the next node does not exist and is not the final node, a new BasicNode will be appended
func (b *BasicNode) append(path string, node PhantomNode) {
	if path == "" {
		b.children = append(b.children, node)
	} else {
		nextName := strings.Split(path, separator)[0]
		next := b.Subset(nextName)
		if next == nil {
			next = newBasicNode(nextName)
			b.children = append(b.children, next)
		}
		nextBasic, ok := next.(*BasicNode)
		if ok {
			nextBasic.append(path[len(nextName):], node)
		}
	}
}

type BasicFile struct {
	os.File
	node PhantomNode
}

func (b *BasicFile) Readdir(count int) ([]os.FileInfo, error) {
	childrenAmount := count
	if count == 0 || len(b.node.Children()) < count {
		childrenAmount = len(b.node.Children())
	}
	infos := make([]os.FileInfo, childrenAmount)
	var finalError error
	var err error
	for i := 0; i < childrenAmount; i++ {
		infos[i], err = b.node.Children()[i].File().Stat()
		if err != nil {
			finalError = err
		}
	}
	return infos, finalError
}

func (b *BasicFile) Stat() (os.FileInfo, error) {
	return BasicStat{node: b.node}, nil
}

type BasicStat struct {
	node PhantomNode
}

func (b BasicStat) Name() string {
	return b.node.Name()
}

func (b BasicStat) Size() int64 {
	return 0
}

func (b BasicStat) Mode() os.FileMode {
	panic("implement me")
}

func (b BasicStat) ModTime() time.Time {
	return time.Unix(rand.Int63(), rand.Int63())
}

func (b BasicStat) IsDir() bool {
	return true
}

func (b BasicStat) Sys() interface{} {
	panic("implement me")
}

func initializeTree() {
	tree = &BasicNode{
		children: make([]PhantomNode, 0),
		name:     "",
	}
	paths := context.Conf.DavPaths
	appendFullPath(paths.PayedBills, newBillCollectionNode(true))

	appendFullPath(paths.UnpayedBills, newBillCollectionNode(false))

	rootCategory := newCategoryNode(&database.Category{Name: paths.Scores}, "")
	rootCategory.root = true
	appendFullPath(paths.Scores, rootCategory)
}

func appendFullPath(path string, node PhantomNode) {
	splitPath := strings.Split(path, separator)
	tree.append(strings.Join(splitPath[:len(splitPath)-1], separator), node)
}

func newBasicNode(name string) *BasicNode {
	node := BasicNode{
		children: make([]PhantomNode, 0),
		name:     name,
	}
	return &node
}
