.PHONY: build

SRC = test.md
PDF = test.pdf
STYLE = breezedark

build:
	@echo "[*] compiling report"
	@pandoc $(SRC) -o $(PDF) \
		--from markdown+yaml_metadata_block+raw_html \
		--to latex \
		--pdf-engine xelatex \
		--standalone \
		--template htb \
		--highlight-style $(STYLE)

