package bibtex

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stormvirux/bibtex/ast"
	"github.com/stormvirux/bibtex/asts"
	"strings"
	"testing"
)

func TestNew_resolve(t *testing.T) {
	tests := []struct {
		src  string
		want Entry
	}{
		{`
				@inproceedings{canonne2020learning,
				  title={Learning from satisfying assignments under continuous distributions},
				  author={Canonne, Clement L and De, Anindya and Servedio, Rocco A},
				  booktitle={Proceedings of the Fourteenth Annual ACM-SIAM Symposium on Discrete Algorithms},
				  pages={82--101},
				  year={2020},
				  organization={SIAM}
			  }
	  `,
			Entry{
				Type: EntryInProceedings,
				Key:  "canonne2020learning",
				Tags: map[Field]ast.Expr{
					"booktitle":    asts.Text("Proceedings of the Fourteenth Annual ACM-SIAM Symposium on Discrete Algorithms"),
					"organization": asts.Text("SIAM"),
					"pages":        asts.Text("82--101"),
					"title":        asts.Text("Learning from satisfying assignments under continuous distributions"),
					"author":       newAuthors(newAuthor("Clement L", "Canonne"), newAuthor("Anindya", "De"), newAuthor("Rocco A", "Servedio")),
					"year":         asts.Text("2020"),
				}},
		},
		{
			`@book{citekey, title={Foo \& Bar \$1} }`,
			Entry{Type: EntryBook, Key: "citekey", Tags: map[Field]ast.Expr{"title": asts.Text("Foo & Bar $1")}},
		},
		{
			`@article{citekey, title={formula $e=mc^2$} }`,
			Entry{Type: EntryArticle, Key: "citekey", Tags: map[Field]ast.Expr{"title": asts.Text("formula $e=mc^2$")}},
		},
		{
			`@book{citekey, url={\url{www.example.com}} }`,
			Entry{Type: EntryBook, Key: "citekey", Tags: map[Field]ast.Expr{"url": asts.Text("www.example.com")}},
		},
		{
			`@article{cite_key, url = "http://example.com/foo--bar/~baz/#" }`,
			Entry{Type: EntryArticle, Key: "cite_key", Tags: map[Field]ast.Expr{"url": asts.Text("http://example.com/foo--bar/~baz/#")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			bib := New(
				WithResolvers(
					NewAuthorResolver("author"),
					ResolverFunc(SimplifyEscapedTextResolver),
					NewRenderParsedTextResolver(),
				))
			file, err := bib.Parse(strings.NewReader(tt.src))
			if err != nil {
				t.Fatal(err)
			}
			entries, err := bib.Resolve(file)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 {
				t.Fatalf("expected exactly 1 entry, got %d entries", len(entries))
			}

			if diff := cmp.Diff(tt.want, entries[0]); diff != "" {
				t.Errorf("Read() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
