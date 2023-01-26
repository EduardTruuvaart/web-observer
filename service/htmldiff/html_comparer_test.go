package htmldiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareDocumentSectionCompareIdenticalSections(t *testing.T) {
	// Arrange
	doc1 := `<html><body><span class="product">Out of Stock</span>My content</body></html>`
	doc2 := `<html><body><span class="product">Out of Stock</span>My content</body></html>`

	// Act
	result, err := CompareDocumentSection(doc1, doc2, "body > span.product")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

}

func TestCompareDocumentSectionCompareDifferentSectionsWithDifferntText(t *testing.T) {
	// Arrange
	doc1 := `<html><body><span class="product">Out of Stock</span>My content</body></html>`
	doc2 := `<html><body><span class="product">In Stock</span>My content</body></html>`

	// Act
	result, err := CompareDocumentSection(doc1, doc2, "body > span.product")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Contains(t, result, "text content: Out of Stock != In Stock")
}

func TestCompareDocumentSectionCompareDifferentSectionsWithDifferntClassAndText(t *testing.T) {
	// Arrange
	doc1 := `<html><body><div class="product-title"><span class="product-soldout">Out of Stock</span>My content</div></body></html>`
	doc2 := `<html><body><div class="product-title"><span class="product-instock">In Stock</span>My content</div></body></html>`

	// Act
	result, err := CompareDocumentSection(doc1, doc2, "div.product-title > span")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, "attribute class: product-soldout != product-instock")
	assert.Contains(t, result, "text content: Out of Stock != In Stock")
}
