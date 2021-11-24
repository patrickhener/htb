.PHONY: build

SRC = ++inmd++
PDF = ++outpdf++

build:
	@echo "[*] compiling report"
	@pandoc $(SRC) -o $(PDF) \
		--from markdown+yaml_metadata_block+raw_html \
		--to latex \
		--pdf-engine xelatex \
		--template htb \
		--listings \
		--top-level-division=chapter
