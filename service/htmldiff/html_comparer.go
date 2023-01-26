package htmldiff

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CompareWholeDocument(html1, html2 string) ([]string, error) {
	var differences []string

	doc1, err := goquery.NewDocumentFromReader(strings.NewReader(html1))
	if err != nil {
		return nil, err
	}

	doc2, err := goquery.NewDocumentFromReader(strings.NewReader(html2))
	if err != nil {
		return nil, err
	}

	// Compare the elements, attributes, and text content of both HTML documents
	doc1.Find("*").Each(func(i int, s *goquery.Selection) {
		el1 := s.Get(0)

		// Find the corresponding element in the second HTML document
		el2 := doc2.Find(el1.DataAtom.String()).Eq(i)

		// Compare the attributes
		for _, attr := range el1.Attr {
			if attr2, exists := el2.Attr(attr.Key); !exists || attr.Val != attr2 {
				differences = append(differences, fmt.Sprintf("Attribute %s: %s != %s", attr.Key, attr.Val, attr2))
			}
		}

		// Compare the text content
		if s.Text() != el2.Text() {
			differences = append(differences, fmt.Sprintf("text content: %s != %s", s.Text(), el2.Text()))
		}
	})

	return differences, nil
}

func CompareDocumentSection(html1, html2, cssSelector string) ([]string, error) {
	var differences []string

	doc1, err := goquery.NewDocumentFromReader(strings.NewReader(html1))
	if err != nil {
		return nil, err
	}

	doc2, err := goquery.NewDocumentFromReader(strings.NewReader(html2))
	if err != nil {
		return nil, err
	}

	// Find the section in both HTML documents
	doc1Section := doc1.Find(cssSelector)
	doc2Section := doc2.Find(cssSelector)

	// Compare the attributes
	doc1Section.Each(func(i int, s *goquery.Selection) {
		el1 := s.Get(0)
		el2 := doc2Section.Eq(i)

		for _, attr := range el1.Attr {
			if attr2, exists := el2.Attr(attr.Key); !exists || attr.Val != attr2 {
				differences = append(differences, fmt.Sprintf("attribute %s: %s != %s", attr.Key, attr.Val, attr2))
			}
		}

		// Compare the text content
		if el1.FirstChild.Data != el2.Get(0).FirstChild.Data {
			differences = append(differences, fmt.Sprintf("text content: %s != %s", el1.FirstChild.Data, el2.Get(0).FirstChild.Data))
		}
	})

	return differences, nil
}
