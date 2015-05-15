package suggest

//Forward Index - Maps the document id to the document
type ForwardIndex map[int]string

func NewForwardIndex() *ForwardIndex {
	i := make(ForwardIndex)
	return &i
}

// AddDoc adds a new documents to the forward index
// and associates it with a docID. For example, for
// doc = "software development", docId = 1
func (x *ForwardIndex) AddDoc(docId int, doc string) {
	(*x)[docId] = doc
}

func (x *ForwardIndex) DocByID(i int) string {
	return (*x)[i]
}
