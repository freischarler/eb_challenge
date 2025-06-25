package mockImpls

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockBooksRepositories_GetBooksProvider(t *testing.T) {
	// Arrange
	mock := NewMockBooksRepositories()
	ctx := context.Background()

	// Act
	books, err := mock.GetBooksProvider(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, books, 3)

	// Verify specific book data
	assert.Equal(t, uint(1), books[0].ID)
	assert.Equal(t, "The Go Programming Language", books[0].Name)
	assert.Equal(t, "Alan Donovan", books[0].Author)
	assert.Equal(t, uint(5000), books[0].UnitsSold)
	assert.Equal(t, uint(40), books[0].Price)

	assert.Equal(t, uint(2), books[1].ID)
	assert.Equal(t, "Clean Code", books[1].Name)
	assert.Equal(t, "Robert C. Martin", books[1].Author)

	assert.Equal(t, uint(3), books[2].ID)
	assert.Equal(t, "The Pragmatic Programmer", books[2].Name)
	assert.Equal(t, "Andrew Hunt", books[2].Author)
}

func TestNewMockBooksRepositories(t *testing.T) {
	// Act
	mock := NewMockBooksRepositories()

	// Assert
	assert.NotNil(t, mock)
}
