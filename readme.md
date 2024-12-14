## Introduction

Once again I'm unable to find a suitable sitemap generator, so I created my own package. This one is optimized for memory efficiency and meets most standard requirements.

## Usage example

Module installation

```bash
go get github.com/CaledoniaProject/gopkg-sitemap
```

Module usage

```go
import (
    "github.com/CaledoniaProject/gopkg-sitemap"
)

// your source of sitemap
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
```

Output directory structure

```bash
/var/www/cdn.example.com/sitemap
- sitemap-index.xml
- file-1.gz
- ...
- file-10.gz
```

