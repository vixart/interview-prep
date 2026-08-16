// Топ-K документов для пользователя.
//
// Хранилище: map DocID -> Document под RWMutex (обновление документа —
// просто перезапись по ключу).
//
// GetTopDocuments: скоринг всех документов распараллелен по чанкам между
// воркерами (CPU-bound), каждый воркер держит СВОЮ min-кучу размера K,
// затем частичные топы сливаются. Итог: O(n log K) вместо полной
// сортировки O(n log n), и без единой точки конкуренции на куче.
package main

import (
	"container/heap"
	"fmt"
	"runtime"
	"sort"
	"sync"
)

type Document struct {
	DocID   string
	Content string
}

type User struct{ UserID string }

type Scorer interface {
	GetScore(doc Document, user User) int
}

type scored struct {
	doc   Document
	score int
}

// minHeap размера K: наверху худший из текущего топа.
type minHeap []scored

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].score < h[j].score }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(scored)) }
func (h *minHeap) Pop() any          { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }

type DocumentService struct {
	mu     sync.RWMutex
	docs   map[string]Document
	scorer Scorer
}

func NewDocumentService(s Scorer) *DocumentService {
	return &DocumentService{docs: make(map[string]Document), scorer: s}
}

// AddDocument добавляет или обновляет документ.
func (s *DocumentService) AddDocument(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.docs[doc.DocID] = doc
	return nil
}

// GetTopDocuments возвращает топ K документов по убыванию скора.
func (s *DocumentService) GetTopDocuments(user User, limit int) ([]Document, error) {
	// Снимок документов под RLock, скоринг — уже без блокировки.
	s.mu.RLock()
	docs := make([]Document, 0, len(s.docs))
	for _, d := range s.docs {
		docs = append(docs, d)
	}
	s.mu.RUnlock()

	workers := runtime.NumCPU()
	chunk := (len(docs) + workers - 1) / workers
	partial := make([]minHeap, workers)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * chunk
		if lo >= len(docs) {
			break
		}
		hi := min(lo+chunk, len(docs))

		wg.Add(1)
		go func(w int, part []Document) {
			defer wg.Done()

			h := make(minHeap, 0, limit)
			for _, d := range part {
				sc := scored{doc: d, score: s.scorer.GetScore(d, user)}
				if len(h) < limit {
					heap.Push(&h, sc)
				} else if sc.score > h[0].score {
					h[0] = sc
					heap.Fix(&h, 0)
				}
			}
			partial[w] = h
		}(w, docs[lo:hi])
	}
	wg.Wait()

	// Слить частичные топы и отсортировать по убыванию скора.
	var all []scored
	for _, h := range partial {
		all = append(all, h...)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].score > all[j].score })

	if len(all) > limit {
		all = all[:limit]
	}
	result := make([]Document, len(all))
	for i, sc := range all {
		result[i] = sc.doc
	}

	return result, nil
}

// lenScorer: скор = длина контента (для демонстрации).
type lenScorer struct{}

func (lenScorer) GetScore(doc Document, _ User) int { return len(doc.Content) }

func main() {
	svc := NewDocumentService(lenScorer{})

	_ = svc.AddDocument(Document{DocID: "doc1", Content: "aaaaaaaaaa"}) // 10
	_ = svc.AddDocument(Document{DocID: "doc2", Content: "aaaaaaaa"})   // 8
	_ = svc.AddDocument(Document{DocID: "doc3", Content: "aaaaaaa"})    // 7
	_ = svc.AddDocument(Document{DocID: "doc4", Content: "a"})          // 1

	top, _ := svc.GetTopDocuments(User{UserID: "user1"}, 2)
	for _, d := range top {
		fmt.Println(d.DocID) // doc1, doc2
	}
}
