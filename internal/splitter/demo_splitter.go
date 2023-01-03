package splitter

func init() {
	//*splitters = append(*splitters, demoSplitterIns)
}

var demoSplitterIns demoSplitter

type demoSplitter struct {
	simpleSplitter
}

func (j demoSplitter) Match(noteType string) bool {
	return noteType == "demo"
}
