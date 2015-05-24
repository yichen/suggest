package suggest

import (
	"strings"
)

type Document struct {
	DocID int
	Bloom int
}

//Inverted Index - Maps the query prefix to the matching documents
type InvertedIndex map[string][]Document

func NewInvertedIndex() *InvertedIndex {
	i := make(InvertedIndex)
	return &i
}

func (x *InvertedIndex) Size() int {
	return len(map[string][]Document(*x))
}

// AddDoc generates the inverted index for the doc string.
// For example, if the input doc is "techcrunch website",
// and the bloom of this doc, it parse the doc into words
// and for each word, take the prefix up to 4 characters, and
// generate map from the prefix to a tuple of docID and the bloom.
// as a result,
// 		t    --> {docId, bloom}
// 		te   --> {docId, bloom}
// 		tec  --> {docId, bloom}
// 		tech --> {docId, bloom}
// 		w 	 --> {docId, bloom}
// 		we   --> {docId, bloom}
// 		web  --> {docId, bloom}
// 		webs --> {docId, bloom}
func (x *InvertedIndex) AddDoc(docId int, doc string, bloom int) {
	for _, word := range strings.Fields(doc) {
		word = getPrefix(word)

		for i := 1; i <= len(word); i++ {
			prefix := word[:i]
			ref, ok := (*x)[prefix]
			if !ok {
				ref = nil
			}

			(*x)[prefix] = append(ref, Document{DocID: docId, Bloom: bloom})
		}
	}
}

func (x *InvertedIndex) Search(query string) []Document {
	q := getPrefix(query)

	ref, ok := (*x)[q]

	if ok {
		return ref
	}
	return nil
}

func getPrefix(query string) string {
	qLen := Min(len(query), 8)
	q := query[0:qLen]
	return strings.ToLower(q)
}

func Min(a ...int) int {
	min := int(^uint(0) >> 1) // largest int
	for _, i := range a {
		if i < min {
			min = i
		}
	}
	return min
}

func Max(a ...int) int {
	max := int(0)
	for _, i := range a {
		if i > max {
			max = i
		}
	}
	return max
}
