package suggest

import (
	"strings"
	"sync"
)

//Used for the bloom filter
const (
	FNV_BASIS_64 = uint64(14695981039346656037)
	FNV_PRIME_64 = uint64((1 << 40) + 435)
	FNV_MASK_64  = uint64(^uint64(0) >> 1)
	NUM_BITS     = 64

	FNV_BASIS_32 = uint32(0x811c9dc5)
	FNV_PRIME_32 = uint32((1 << 24) + 403)
	FNV_MASK_32  = uint32(^uint32(0) >> 1)
)

// Suggest is the container for the indices and provide the
// entry point for main public methods
type Suggest struct {
	//maxDocId is to keep track of the last docId so each
	//new doc will have a unique id
	nextID int
	iIndex *InvertedIndex
	fIndex *ForwardIndex

	docs map[string]bool

	sync.Mutex
}

func NewSuggest() *Suggest {
	suggest := &Suggest{
		nextID: 0,
		iIndex: NewInvertedIndex(),
		fIndex: NewForwardIndex(),
		docs:   make(map[string]bool),
	}

	return suggest
}

func (s *Suggest) AddDocument(doc string) {

	s.Lock()
	defer s.Unlock()

	filter := computeBloomFilter(doc)
	s.iIndex.AddDoc(s.nextID, doc, filter)
	s.fIndex.AddDoc(s.nextID, doc)
	s.docs[doc] = true
	s.nextID++
}

func (s *Suggest) ContainsDocument(doc string) bool {
	if s.docs[doc] == true {
		return true
	}
	return false
}

func (s *Suggest) AddSymbol(symbol string) {
	s.Lock()
	defer s.Unlock()

	doc := tokenizeSymbol(symbol)
	filter := computeBloomFilter(doc)
	s.iIndex.AddDoc(s.nextID, doc, filter)

	// the forward index will retrieve the original symbol
	// instead of the tokenized doc
	s.fIndex.AddDoc(s.nextID, symbol)
	s.docs[symbol] = true
	s.nextID++
}

func (s *Suggest) Search(query string) []string {

	// first search the inverted index for some candidate
	candidates := s.iIndex.Search(query)

	// filter out the candidate using bloom filter
	queryBloom := computeBloomFilter(query)

	// docIDs is the set of result (to remove duplicates)
	docIDs := make(map[int]bool)
	for _, i := range candidates {
		if TestBytesFromQuery(i.Bloom, queryBloom) == true {
			docIDs[i.DocID] = true
		}
	}

	result := make([]string, 0, 0)
	for docID, _ := range docIDs {
		result = append(result, s.fIndex.DocByID(docID))
	}

	SortByRank(query, result)

	return result
}

func (s *Suggest) SearchAll(query string) []string {
	resultSet := map[string]int{}

	// search for each word in the query string
	// and record the number of times in each resultset
	words := strings.Fields(query)
	for _, word := range words {
		result := s.Search(word)

		for _, r := range result {
			if c, ok := resultSet[r]; ok {
				resultSet[r] = c + 1
			} else {
				resultSet[r] = 1
			}
		}
	}

	// The final result is by finding the intersection
	// of all the result set
	finalResultSet := []string{}
	for doc, count := range resultSet {
		if count == len(words) {
			finalResultSet = append(finalResultSet, doc)
		}
	}

	SortByRank(query, finalResultSet)

	return finalResultSet
}

//The bloom filter of a word is 8 bytes in length
//and has each character added separately. That is,
// if the word is "software development", we make sure
// each of the character in this word, [s, o, f, t, w, a, r, e, d, v, l, p, m, n],
// will be mapped to bits in the filter. This bloom is used to test if
// a given input string, like "softdev", that has all its characters included
// in the document
func computeBloomFilter(s string) int {
	cnt := len(s)

	if cnt <= 0 {
		return 0
	}

	var filter int
	hash := uint64(0)

	for i := 0; i < cnt; i++ {
		c := s[i]

		//first hash function
		hash ^= uint64(0xFF & c)
		hash *= FNV_PRIME_64

		//second hash function (reduces collisions for bloom)
		hash ^= uint64(0xFF & (c >> 16))
		hash *= FNV_PRIME_64

		//position of the bit mod the number of bits (8 bytes = 64 bits)
		bitpos := hash % NUM_BITS
		if bitpos < 0 {
			bitpos += NUM_BITS
		}
		filter = filter | (1 << bitpos)
	}

	return filter
}

//Iterates through all of the 8 bytes (64 bits) and tests
//each bit that is set to 1 in the query's filter against
//the bit in the comparison's filter.  If the bit is not
// also 1, you do not have a match.
func TestBytesFromQuery(bf int, qBloom int) bool {
	return (bf & qBloom) != 0
}

func tokenizeSymbol(symbol string) string {
	sanitized := ""

	isLowerCase := true

	for _, runeValue := range symbol {
		if runeValue >= 'a' && runeValue <= 'z' || runeValue >= 'A' && runeValue <= 'Z' || runeValue >= '0' && runeValue <= '9' {
			if runeValue > 'Z' || runeValue < 'A' {
				isLowerCase = true
			}
			// if transition from lowercase to uppercase, split the word
			if runeValue >= 'A' && runeValue <= 'Z' && isLowerCase {
				isLowerCase = false
				sanitized += " "
			}
			sanitized += string(runeValue)
		} else {
			sanitized += " "
		}
	}

	return strings.ToLower(strings.Trim(sanitized, " "))
}
