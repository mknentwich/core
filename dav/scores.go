package dav

import "golang.org/x/net/webdav"

type ScoreCollectionNode struct{}

func (s *ScoreCollectionNode) File() webdav.File {
	panic("implement me")
}

func (s *ScoreCollectionNode) Name() string {
	panic("implement me")
}

func (s *ScoreCollectionNode) Subset(string) PhantomNode {
	panic("implement me")
}

func (s *ScoreCollectionNode) append(string, PhantomNode) {
	panic("implement me")
}
