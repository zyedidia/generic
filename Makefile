DOCS=avl/README.md btree/README.md cache/README.md hashmap/README.md hashset/README.md interval/README.md iter/README.md list/README.md rope/README.md stack/README.md DOC.md

all: $(DOCS)

%/README.md: %
	gomarkdoc --output $@ ./$<

DOC.md: .
	gomarkdoc --output $@ .
