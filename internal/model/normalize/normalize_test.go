package normalize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Category(t *testing.T) {
	t.Run("обычное слово, кириллица, ловеркейс", func(t *testing.T) {
		res := Category("слово")
		assert.Equal(t, "Слово", res)
	})

	t.Run("обычное слово, кириллица, капс", func(t *testing.T) {
		res := Category("СЛОВО")
		assert.Equal(t, "Слово", res)
	})

	t.Run("обычное слово, кириллица, смешанный", func(t *testing.T) {
		res := Category("СлОво")
		assert.Equal(t, "Слово", res)
	})

	t.Run("слово с пробелом и дефисом, кириллица, смешанный", func(t *testing.T) {
		res := Category("какОЕ-то СлОво")
		assert.Equal(t, "Какое-то слово", res)
	})

	t.Run("слово с дефисом и пробелом, латиница, смешанный", func(t *testing.T) {
		res := Category("sOme w-OrD")
		assert.Equal(t, "Some w-ord", res)
	})

	t.Run("пустая строка", func(t *testing.T) {
		res := Category("")
		assert.Equal(t, "", res)
	})
}
