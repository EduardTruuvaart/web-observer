package htmldiff

import (
	"fmt"
	"strings"

	domain "github.com/EduardTruuvaart/web-observer/domain/htmlcompare"
	"github.com/PuerkitoBio/goquery"
)

func CompareDocumentSection(sorceHtml, targetHtml, cssSelector string) (domain.HtmlCompareResult, error) {
	var differences []string
	var result domain.HtmlCompareResult

	doc1, err := goquery.NewDocumentFromReader(strings.NewReader(sorceHtml))
	if err != nil {
		return result, err
	}

	doc2, err := goquery.NewDocumentFromReader(strings.NewReader(targetHtml))
	if err != nil {
		return result, err
	}

	// Find the section in both HTML documents
	doc1Section := doc1.Find(cssSelector)
	doc2Section := doc2.Find(cssSelector)

	if doc1Section.Length() == 0 {
		return domain.HtmlCompareResult{
			State: domain.SelectionNotFoundInSource,
		}, nil
	}

	if doc2Section.Length() == 0 {
		return domain.HtmlCompareResult{
			State: domain.SelectionNotFoundInTarget,
		}, nil
	}

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

	if len(differences) == 0 {
		return domain.HtmlCompareResult{
			State: domain.Identical,
		}, nil
	}

	return domain.HtmlCompareResult{
		State:       domain.Different,
		DiffSize:    len(differences),
		Differences: differences,
	}, nil
}
