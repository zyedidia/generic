DOCS=avl/README.md btree/README.md cache/README.md hashmap/README.md hashset/README.md interval/README.md list/README.md mapset/README.md multimap/README.md rope/README.md stack/README.md trie/README.md DOC.md queue/README.md heap/README.md

all: $(DOCS)

%/README.md: %
	gomarkdoc --output $@ ./$<

DOC.md: .
	gomarkdoc --output $@ .
