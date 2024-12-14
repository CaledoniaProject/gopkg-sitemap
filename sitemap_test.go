package sitemap_test

import (
	"testing"
	"time"

	sitemap "github.com/CaledoniaProject/gopkg-sitemap"
)

type MyArticle struct {
	ThumbnailURL []string
	CanonicalURL string
	CreatedAt    time.Time
}

// major function to create sitemap and index file
func WriteSitemap(articles []*MyArticle) error {
	sg := &sitemap.SitemapGenerator{
		OutputDirectory: "/tmp/sitemap", // where to save these files
		LinksPerFile:    5000,           // adjust this parameter to avoid search engine limits
	}

	for _, article := range articles {
		sitemapURL := &sitemap.SitemapURL{
			Loc:        article.CanonicalURL,
			LastMod:    article.CreatedAt,
			ChangeFreq: sitemap.Daily,
			Priority:   0.4,
		}

		for _, thumbnail := range article.ThumbnailURL {
			sitemapURL.Images = append(sitemapURL.Images, &sitemap.SitemapImage{Loc: thumbnail})
		}

		if err := sg.AddURL(sitemapURL); err != nil {
			return err
		}
	}

	return sg.WriteIndex("http://cdn.example.com")
}

func TestSitemap(t *testing.T) {
	articles := []*MyArticle{
		{
			ThumbnailURL: []string{"http://cdn.example.com/1.jpg"},
			CanonicalURL: "https://www.example.com/article1",
			CreatedAt:    time.Now(),
		},
	}

	if err := WriteSitemap(articles); err != nil {
		t.Fatalf("faild: %v", err)
	}
}
