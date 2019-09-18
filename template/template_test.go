package template

import (
	"errors"
	"fmt"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	log2 "log"
	"net/http"
	"testing"
)

var categories = []database.Category{
	{Name: "Märsche", Children: []database.Category{{Name: "Behmische Märsche", Scores: []database.Score{{Title: "Koline"}, {Title: "Andulka"}}}}, Scores: []database.Score{{Title: "Arnheim"}}},
	{Name: "Polka", Scores: []database.Score{{Title: "Ma posledni"}}}}

func TestGenerate(t *testing.T) {
	log2.SetFlags(log2.Llongfile | log2.Flags())
	go func() {
		context.InitializeCustomConfig(map[string]context.Serve{
			"db":       database.Serve,
			"template": Serve},
			&context.Configuration{
				GeneratedDirectory: "gen",
				Host:               "0.0.0.0:9400",
				SQLiteFile:         ":memory:",
			})
	}()
	for err := errors.New(""); err != nil; _, err = http.Get("http://127.0.0.1:9400/") {
	}
	rest.InsertNewCategory(categories[0])
	rest.InsertNewCategory(categories[1])
	if err := Generate(); err != nil {
		fmt.Println(err.Error())
	}
}
