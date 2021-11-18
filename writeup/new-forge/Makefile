.PHONY: build

SRC = forge-writeup.md
PDF = forge-writeup.pdf

build:
	@echo "[*] compiling report"
	@pandoc $(SRC) -o $(PDF) \
		--from markdown+yaml_metadata_block+raw_html \
		--to latex \
		--pdf-engine xelatex \
		--standalone \
		--template htb \
		--listings