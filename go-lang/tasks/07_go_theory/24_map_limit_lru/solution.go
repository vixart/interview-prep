// Ограничение размера мапы — LRU-вытеснение.
//
// В самой мапе порядок вставки не хранится, поэтому "самую старую" запись
// на голой мапе не найти. Классическое решение — LRU: двусвязный список
// (container/list) хранит порядок использования, мапа даёт O(1) доступ
// к узлам списка. При превышении лимита удаляем элемент с хвоста.
package main

import (
	"container/list"
	"fmt"
)

type WordCounter struct {
	counts map[string]*list.Element
	order  *list.List // голова — самые свежие, хвост — самые старые
	limit  int
}

type entry struct {
	word  string
	count int
}

func NewWordCounter(limit int) *WordCounter {
	return &WordCounter{
		counts: make(map[string]*list.Element, limit),
		order:  list.New(),
		limit:  limit,
	}
}

func (wc *WordCounter) CountWord(word string) {
	if el, ok := wc.counts[word]; ok {
		el.Value.(*entry).count++
		wc.order.MoveToFront(el) // слово снова использовано — освежаем
		return
	}

	wc.counts[word] = wc.order.PushFront(&entry{word: word, count: 1})

	if wc.order.Len() > wc.limit {
		oldest := wc.order.Back()
		wc.order.Remove(oldest)
		delete(wc.counts, oldest.Value.(*entry).word)
	}
}

func (wc *WordCounter) Counts() map[string]int {
	out := make(map[string]int, len(wc.counts))
	for w, el := range wc.counts {
		out[w] = el.Value.(*entry).count
	}
	return out
}

func main() {
	wc := NewWordCounter(3)

	words := []string{"apple", "banana", "apple", "orange", "grape", "banana", "kiwi"}
	for _, w := range words {
		wc.CountWord(w)
	}

	// В мапе не больше 3 слов; вытеснены самые давно использованные.
	fmt.Println("Количество слов:", wc.Counts())
	// map[banana:2 grape:1 kiwi:1]
}
