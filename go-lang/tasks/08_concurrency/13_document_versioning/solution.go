// Версионирование документов из очереди: произвольный порядок, дубликаты,
// конкурентный доступ.
//
// Состояние по каждому Url — актуальный документ + множество виденных
// FetchTime (для отсечения дубликатов). Правила слияния:
//   - пришла более свежая версия -> обновляем Text/FetchTime;
//   - пришла более старая        -> обновляем только PubDate/FirstFetchTime;
//   - FirstFetchTime/PubDate всегда от МИНИМАЛЬНОГО FetchTime.
//
// Наружу отдаём копию, чтобы вызывающий не мог менять внутреннее состояние.
package main

import (
	"fmt"
	"sync"
)

type Document struct {
	Url            string
	PubDate        uint64
	FetchTime      uint64
	Text           string
	FirstFetchTime *uint64
}

type state struct {
	doc  Document
	seen map[uint64]struct{} // виденные FetchTime — фильтр дубликатов
}

type Processor struct {
	mu   sync.Mutex
	docs map[string]*state
}

func NewProcessor() *Processor {
	return &Processor{docs: make(map[string]*state)}
}

// Process возвращает актуальную версию или nil для дубликата.
func (p *Processor) Process(doc Document) (*Document, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	st, ok := p.docs[doc.Url]
	if !ok {
		first := doc.FetchTime
		doc.FirstFetchTime = &first
		p.docs[doc.Url] = &state{
			doc:  doc,
			seen: map[uint64]struct{}{doc.FetchTime: {}},
		}
		return copyDoc(doc), nil
	}

	if _, dup := st.seen[doc.FetchTime]; dup {
		return nil, nil // дубликат — обновление не требуется
	}
	st.seen[doc.FetchTime] = struct{}{}

	cur := &st.doc
	if doc.FetchTime > cur.FetchTime {
		// более свежая версия
		cur.FetchTime = doc.FetchTime
		cur.Text = doc.Text
	}
	if doc.FetchTime < *cur.FirstFetchTime {
		// более старая версия — она и есть «первая»
		*cur.FirstFetchTime = doc.FetchTime
		cur.PubDate = doc.PubDate
	}

	return copyDoc(*cur), nil
}

func copyDoc(d Document) *Document {
	c := d
	if d.FirstFetchTime != nil {
		v := *d.FirstFetchTime
		c.FirstFetchTime = &v
	}
	return &c
}

func show(d *Document) string {
	if d == nil {
		return "<dup>"
	}
	return fmt.Sprintf("{%s fetch=%d text=%s pub=%d first=%d}",
		d.Url, d.FetchTime, d.Text, d.PubDate, *d.FirstFetchTime)
}

func main() {
	p := NewProcessor()

	d1, _ := p.Process(Document{Url: "doc1", FetchTime: 100, Text: "v1", PubDate: 50})
	fmt.Println(show(d1)) // {doc1 fetch=100 text=v1 pub=50 first=100}

	d2, _ := p.Process(Document{Url: "doc1", FetchTime: 200, Text: "v2", PubDate: 60})
	fmt.Println(show(d2)) // {doc1 fetch=200 text=v2 pub=50 first=100}

	d3, _ := p.Process(Document{Url: "doc1", FetchTime: 50, Text: "v0", PubDate: 40})
	fmt.Println(show(d3)) // {doc1 fetch=200 text=v2 pub=40 first=50}

	dup, _ := p.Process(Document{Url: "doc1", FetchTime: 200, Text: "v2", PubDate: 60})
	fmt.Println(show(dup)) // <dup>

	// Конкурентная обработка: состояние защищено мьютексом.
	var wg sync.WaitGroup
	for i := uint64(1); i <= 100; i++ {
		wg.Add(1)
		go func(i uint64) {
			defer wg.Done()
			_, _ = p.Process(Document{Url: "doc2", FetchTime: i, Text: "t", PubDate: i})
		}(i)
	}
	wg.Wait()

	final, _ := p.Process(Document{Url: "doc2", FetchTime: 101, Text: "last", PubDate: 1})
	fmt.Println(show(final)) // {doc2 fetch=101 text=last pub=1 first=1}
}
