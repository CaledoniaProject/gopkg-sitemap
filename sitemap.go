package sitemap

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

const (
	urlSetStart = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1">`
	urlSetEnd   = `</urlset>`
	indexStart  = `<?xml version="1.0" encoding="UTF-8"?><sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexEnd    = `</sitemapindex>`
)

type ChangeFreq string

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

type SitemapImage struct {
	Loc string `xml:"image:loc"`
}

type SitemapURL struct {
	XMLName    xml.Name        `xml:"url"`
	Loc        string          `xml:"loc"`
	LastMod    time.Time       `xml:"lastmod"`
	ChangeFreq ChangeFreq      `xml:"changefreq"`
	Priority   float64         `xml:"priority"`
	Images     []*SitemapImage `xml:"image:image"`
}

type SitemapLoc struct {
	XMLName xml.Name  `xml:"sitemap"`
	Loc     string    `xml:"loc"`
	LastMod time.Time `xml:"lastmod"`
}

type SitemapGenerator struct {
	// configuration
	OutputDirectory string
	LinksPerFile    int

	// internal variables
	Files         []string
	writer        *gzip.Writer
	numberOfLinks int
}

func (s *SitemapGenerator) AddURL(sitemapURL *SitemapURL) error {
	// create writer
	if s.writer == nil {
		filename := fmt.Sprintf("file-%d.gz", len(s.Files)+1)

		// create directory
		if err := os.MkdirAll(s.OutputDirectory, 0755); err != nil {
			return err
		} else if filp, err := os.OpenFile(s.OutputDirectory+"/"+filename, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		} else {
			s.writer = gzip.NewWriter(filp)
			s.Files = append(s.Files, filename)
		}

		if _, err := s.writer.Write([]byte(urlSetStart)); err != nil {
			return err
		}
	}

	// write marshalled xml
	if xmlData, err := xml.Marshal(sitemapURL); err != nil {
		return err
	} else if _, err := s.writer.Write(xmlData); err != nil {
		return err
	} else {
		s.numberOfLinks++

		// move to next file
		if s.numberOfLinks == s.LinksPerFile {
			s.Close()
		}
	}

	return nil
}

func (s *SitemapGenerator) Close() error {
	// no active file
	if s.writer == nil {
		return nil
	}

	// write closing header and close
	if _, err := s.writer.Write([]byte(urlSetEnd)); err != nil {
		return err
	} else if err := s.writer.Close(); err != nil {
		return err
	}

	// reset staths
	s.numberOfLinks = 0
	s.writer = nil
	return nil
}

func (s *SitemapGenerator) WriteIndex(baseURL string) error {
	if err := s.Close(); err != nil {
		return err
	}

	filp, err := os.OpenFile(s.OutputDirectory+"/sitemap-index.xml", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	} else if _, err := filp.Write([]byte(indexStart)); err != nil {
		return err
	}

	for _, file := range s.Files {
		if data, err := xml.Marshal(&SitemapLoc{
			Loc:     baseURL + "/" + file,
			LastMod: time.Now(),
		}); err != nil {
			return err
		} else if _, err := filp.Write(data); err != nil {
			return err
		}
	}

	if _, err := filp.Write([]byte(indexEnd)); err != nil {
		return err
	}

	return filp.Close()
}
