package dav

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"golang.org/x/net/webdav"
	"strings"
)

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

type CategoryNode struct {
	*BasicNode
	parentPath string
}

func (c *CategoryNode) Path() string {
	return c.parentPath + separator + c.Name()
}

//path is the name to this node, i have to check only it's children
func (c *CategoryNode) Children() []PhantomNode {
	nodes := make([]PhantomNode, 0)
	category := categoryAt(c.Path(), rest.QueryCategoriesWithChildrenAndScoresPreserve())
	nodes = make([]PhantomNode, len(category.Children)+len(category.Scores))
	i := 0
	for _, child := range category.Children {
		nodes[i] = newCategoryNode(&child, c.Path())
		i++
	}
	for _, score := range category.Scores {
		nodes[i] = newScoreNode(&score)
		i++
	}
	return nodes
}

func categoryAt(path string, categories []database.Category) database.Category {
	names := strings.Split(path, separator)
	for _, category := range categories {
		if category.Name == names[0] {
			if len(names) == 1 {
				return category
			} else {
				return categoryAt(strings.Join(names[1:], separator), category.Children)
			}
		}
	}
	log(context.LOG_WARNING, "tried to fetch a category which does not exist: %s", path)
	return database.Category{Name: "Black Hole"}
}

func newCategoryNode(category *database.Category, parentPath string) *CategoryNode {
	return &CategoryNode{
		BasicNode: &BasicNode{
			name: category.Name},
		parentPath: parentPath,
	}
}

func newScoreNode(score *database.Score) *ScoreNode {
	return &ScoreNode{
		BasicNode: &BasicNode{
			name: score.Title},
	}
}

type ScoreNode struct {
	*BasicNode
}
