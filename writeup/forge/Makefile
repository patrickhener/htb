# https://github.com/shiblon/latex-makefile

MAKEFILE_VERSION := 1.0.0-0
.DEFAULT_GOAL	:= all
# Note that the user-global version is imported *after* the source directory,
# so that you can use stuff like ?= to get proper override behavior.
.PHONY: Makefile GNUmakefile Makefile.ini $(HOME)/.latex-makefile/Makefile.ini
-include Makefile.ini
-include $(HOME)/.latex-makefile/Makefile.ini

# MacOS fixes
#############
ifneq (,$(findstring sha1sum,$(shell which sha1sum)))
	SHA1SUM := sha1sum
endif
ifneq (,$(findstring shasum,$(shell which shasum)))
	SHA1SUM := shasum
endif

SED ?= LC_ALL=C sed
#############

TMP_DIR := ${XDG_RUNTIME_DIR}
ifeq ($(TMP_DIR),)
	TMP_DIR := /tmp
endif
OUTPUT_DIR := $(TMP_DIR)/latex/$(shell pwd -P|$(SHA1SUM)|cut -d" " -f1|head -c8)
$(shell mkdir -p '$(OUTPUT_DIR)')
$(shell chmod 700 '$(OUTPUT_DIR)')

SCRIPT_DIR := $(shell dirname -- "`kpsewhich writeup.cls`")/scripts/

BINARY_TARGET_DIR ?= .
ADDITIONAL_LATEX_ARGS ?=
LATEX_ARGS ?= $(ADDITIONAL_LATEX_ARGS) -output-directory='$(OUTPUT_DIR)' -halt-on-error
CLASSOPTIONS_REPORT ?=

latex_build_program ?= lualatex
build_target_extension ?= pdf
IGNORE_TARGETS ?=
N ?= 5

# Output format
max_print_line=1000
error_line=254
half_error_line=238
export max_print_line
export error_line
export half_error_line

# Files of interest
all_files.tex		?= $(wildcard *.tex)

# Utility function for getting all .$1 files that are to be ignored
#  * files listed in $(includes.$1)
#  * files not listed in $(onlysources.$1) if it is defined
ignore_files = \
  $(includes.$1) \
  $(if $(onlysources.$1),$(filter-out $(onlysources.$1), $(all_files.$1)))

# Patterns to never be allowed as source targets
ignore_patterns	:= %._include_

# Patterns allowed as source targets but not included in 'all' builds
nodefault_patterns := %._nobuild_ $(ignore_patterns)

# Utility function for getting targets suitable building
# $(call filter-buildable,suffix)
filter-buildable	= \
	$(filter-out $(call ignore_files,$1) \
		$(addsuffix .$1,$(ignore_patterns)),$(all_files.$1))

# Utility function for getting targets suitable for 'all' builds
# $(call filter-default,suffix)
filter-default		= \
	$(filter-out $(call ignore_files,$1) \
		$(addsuffix .$1,$(nodefault_patterns)),$(all_files.$1))

# Top level sources that can be built even when they are not by default
files.tex	:= $(call filter-buildable,tex)

# Utility function for obtaining stems
# $(call get-stems,suffix,[prefix])
get-stems	= $(sort $($(if $2,$2_,)files.$1:%.$1=%))

# Top level sources that are built by default targets
default_files.tex	:= $(call filter-default,tex)

# List of all default stems (all default PDF targets):
default_stems.tex		:= $(call get-stems,tex,default)

# List of all stems (all possible bare PDF targets created here):
stems.tex		:= $(call get-stems,tex)

# Utility function for creating larger lists of stems
# $(call concat-stems,suffixes,[prefix])
concat-stems	= $(sort $(foreach s,$1,$($(if $2,$2_,)stems.$s)))

stems_source		:= $(call concat-stems,tex)
stems_script		:= $(call concat-stems)
stems_ss		:= $(sort $(stems_source) $(stems_script))
stems_ssg		:= $(sort $(stems_ss))
all_pdf_targets		:= $(addsuffix .pdf,$(stems_ssg))

rename-report = '$(SCRIPT_DIR)rename-report.sh'

.PHONY: all
all: $(all_pdf_targets)

.PHONY: cleanbuilddir
cleanbuilddir:
	@rm -rf "$(OUTPUT_DIR)"/* 2> /dev/null

.PHONY: clean
clean: cleanbuilddir
	@rm -f $(foreach stem,$(stems.tex),$(addprefix $(stem),.pdf .csv .log .aux .toc .d .fls)) 2> /dev/null
	@rm -f ./*.csv 2> /dev/null
	@rm -f ./*.pdf 2> /dev/null


tex-argument = "$(addprefix \PassOptionsToClass{,$(addsuffix }{writeup},$(CLASSOPTIONS_REPORT)))$(addprefix \PassOptionsToClass{,$(addsuffix }{wu_main},$(CLASSOPTIONS_MAIN)))\input{./$1}"

# LaTeX invocations
#
# $(call latex,<tex file>,[<extra LaTeX args>])
run-latex	= [ -n '$(tmp_texmfhome)' ] && export TEXMFHOME='$(tmp_texmfhome)'; $(latex_build_program) $(LATEX_ARGS) -interaction=batchmode -file-line-error $(if $2,$2,) $(call tex-argument,$1) > /dev/null

# Don't call this directly - it is here to avoid calling wildcard more than
# once in remove-files.
remove-files-helper	= $(if $1,rm $1,$(sh_true))

# $(call remove-files,file1 file2)
remove-files		= $(call remove-files-helper,$(wildcard $1))

# Terminal color definitions

# If NO_COLOR is not explicitly set, set it to true if stdin is not a TTY.
# Note: Testing stdout instead of stdin seems more sensible, but in fact does not work, since the shell's stdout is always make.
ifndef NO_COLOR
    NO_COLOR := $(shell [ -t 0 ] || echo true)
endif
REAL_TPUT 	:= $(if $(NO_COLOR),,$(shell which tput))

# $(call get-term-code,codeinfo)
# e.g.,
# $(call get-term-code,setaf 0)
get-term-code = $(if $(REAL_TPUT),$(shell $(REAL_TPUT) $1),)

black	:= $(call get-term-code,setaf 0)
red	:= $(call get-term-code,setaf 1)
green	:= $(call get-term-code,setaf 2)
yellow	:= $(call get-term-code,setaf 3)
blue	:= $(call get-term-code,setaf 4)
magenta	:= $(call get-term-code,setaf 5)
cyan	:= $(call get-term-code,setaf 6)
white	:= $(call get-term-code,setaf 7)
bold	:= $(call get-term-code,bold)
uline	:= $(call get-term-code,smul)
reset	:= $(call get-term-code,sgr0)

#
# User-settable definitions
#
LATEX_COLOR_WARNING	?= magenta
LATEX_COLOR_ERROR	?= red
LATEX_COLOR_INFO	?= green
LATEX_COLOR_UNDERFULL	?= magenta
LATEX_COLOR_OVERFULL	?= red bold
LATEX_COLOR_UREFERENCES	?= red bold
LATEX_COLOR_PAGES	?= bold
LATEX_COLOR_BUILD	?= cyan
LATEX_COLOR_GRAPHIC	?= yellow
LATEX_COLOR_DEP		?= green
LATEX_COLOR_SUCCESS	?= green bold
LATEX_COLOR_FAILURE	?= red bold
# Gets the real color from a simple textual definition like those above
# $(call get-color,ALL_CAPS_COLOR_NAME)
# e.g., $(call get-color,WARNING)
get-color	= $(subst $(space),,$(foreach c,$(LATEX_COLOR_$1),$($c)))

#
# STANDARD COLORS
#
C_WARNING	:= $(call get-color,WARNING)
C_ERROR		:= $(call get-color,ERROR)
C_INFO		:= $(call get-color,INFO)
C_UNDERFULL	:= $(call get-color,UNDERFULL)
C_OVERFULL	:= $(call get-color,OVERFULL)
C_UREFERENCES	:= $(call get-color,UREFERENCES)
C_PAGES		:= $(call get-color,PAGES)
C_BUILD		:= $(call get-color,BUILD)
C_GRAPHIC	:= $(call get-color,GRAPHIC)
C_DEP		:= $(call get-color,DEP)
C_SUCCESS	:= $(call get-color,SUCCESS)
C_FAILURE	:= $(call get-color,FAILURE)
C_RESET		:= $(reset)

# Colorize LaTeX output.
color_tex := \
$(SED) \
-e '$${' \
-e '  /^$$/!{' \
-e '    H' \
-e '    s/.*//' \
-e '  }' \
-e '}' \
-e '/^$$/!{' \
-e '  H' \
-e '  d' \
-e '}' \
-e '/^$$/{' \
-e '  x' \
-e '  s/^\n//' \
-e '  /Output written on /{' \
-e '    s/.*Output written on \([^(]*\) (\([^)]\{1,\}\)).*/Success!  Wrote \2 to \1/' \
-e '    s/[[:digit:]]\{1,\}/$(C_PAGES)&$(C_RESET)/g' \
-e '    s/Success!/$(C_SUCCESS)&$(C_RESET)/g' \
-e '    s/to \(.*\)$$/to $(C_SUCCESS)\1$(C_RESET)/' \
-e '    b end' \
-e '  }' \
-e '  / *LaTeX Error:.*/{' \
-e '    s/.*\( *LaTeX Error:.*\)/$(C_ERROR)\1$(C_RESET)/' \
-e '    b end' \
-e '  }' \
-e '  /.*LaTeX Warning: There were undefined references.*/{' \
-e '    s/.*\(LaTeX Warning: There were undefined references.*\)/$(C_UREFERENCES)\1$(C_RESET)/' \
-e '    b end' \
-e '  }' \
-e '  /^[[:print:]]* Warning:.*/{' \
-e '    s//$(C_WARNING)&$(C_RESET)/' \
-e '    b end' \
-e '  }' \
-e '  /Underfull.*/{' \
-e '    s/.*\(Underfull.*\)/$(C_UNDERFULL)\1$(C_RESET)/' \
-e '    b end' \
-e '  }' \
-e '  /Overfull.*/{' \
-e '    s/.*\(Overfull.*\)/$(C_OVERFULL)\1$(C_RESET)/' \
-e '    b end' \
-e '  }' \
-e '  d' \
-e '  :end' \
-e '  G' \
-e '}'

# $(call latex-color-log,<LaTeX stem>)
latex-color-log	= $(color_tex) '$1.log'

# Colorizes real, honest-to-goodness LaTeX errors that can't be overcome with
# recompilation.
#
# Note that we only ignore file not found errors for things that we know how to
# build, like graphics files.
#
# Also note that the output of this is piped through sed again to escape any
# backslashes that might have made it through.  This is to avoid sending things
# like "\right" to echo, which interprets \r as LF.  In bash, we could just do
# ${var//\\/\\\\}, but in other popular sh variants (like dash), this doesn't
# work.
#
# $(call colorize-latex-errors,<log file>)
define colorize-latex-errors
$(SED) \
-e '$${' \
-e '  /^$$/!{' \
-e '    H' \
-e '    s/.*//' \
-e '  }' \
-e '}' \
-e '/^$$/!{' \
-e '  H' \
-e '  d' \
-e '}' \
-e '/^$$/{' \
-e '  x' \
-e '  s/^\(\n\)\(.*\)/\2\1/' \
-e '}' \
-e '/^::P\(P\{1,\}\)::/{' \
-e '  s//::\1::/' \
-e '  G' \
-e '  h' \
-e '  d' \
-e '}' \
-e '/^::P::/{' \
-e '  s//::0::/' \
-e '  G' \
-e '}' \
-e 'b start' \
-e ':needonemore' \
-e 's/^/::P::/' \
-e 'G' \
-e 'h' \
-e 'd' \
-e ':needtwomore' \
-e 's/^/::PP::/' \
-e 'G' \
-e 'h' \
-e 'd' \
-e ':needthreemore' \
-e 's/^/::PPP::/' \
-e 'G' \
-e 'h' \
-e 'd' \
-e ':start' \
-e '/^! LaTeX Error: File /{' \
-e '  s/\n//g' \
-e '  b needtwomore' \
-e '}' \
-e 's/^[^[:cntrl:]:]*:[[:digit:]]\{1,\}:/!!! &/' \
-e 's/^\(.*\n\)\([^[:cntrl:]:]*:[[:digit:]]\{1,\}: .*\)/\1!!! \2/' \
-e '/^!!! .* LaTeX Error: File /{' \
-e '  s/\n//g' \
-e '  b needonemore' \
-e '}' \
-e '/^::0::! LaTeX Error: File .*/{' \
-e '  /\n\n$$/{' \
-e '    s/^::0:://' \
-e '    b needonemore' \
-e '  }' \
-e '  s/^::0::! //' \
-e '  s/^\(.*not found.\).*Enter file name:.*\n\(.*[[:digit:]]\{1,\}\): Emergency stop.*/\2: \1/' \
-e '  b error' \
-e '}' \
-e '/^::0::!!! .*LaTeX Error: File .*/{' \
-e '  /\n\n$$/{' \
-e '    s/^::0:://' \
-e '    b needonemore' \
-e '  }' \
-e '  s/::0::!!! //' \
-e '  /could not locate.*any of these extensions:/{' \
-e '    d' \
-e '  }' \
-e '  s/\(not found\.\).*/\1/' \
-e '  b error' \
-e '}' \
-e '/^\(.* LaTeX Error: Missing .begin.document.\.\).*/{' \
-e '  s//\1 --- Are you trying to build an include file?/' \
-e '  b error' \
-e '}' \
-e '/.*\(!!! .*Undefined control sequence\)[^[:cntrl:]]*\(.*\)/{' \
-e '  s//\1: \2/' \
-e '  s/\nl\.[[:digit:]][^[:cntrl:]]*\(\\[^\\[:cntrl:]]*\).*/\1/' \
-e '  b error' \
-e '}' \
-e '/^\(!pdfTeX error:.*\)s*/{' \
-e '  b error' \
-e '}' \
-e '/.*\(!!! .*\)/{' \
-e '  s//\1/' \
-e '  s/[[:cntrl:]]//' \
-e '  s/[[:cntrl:]]$$//' \
-e '  b error' \
-e '}' \
-e 'd' \
-e ':error' \
-e 's/^!\(!! \)\{0,1\}\(.*\)/$(C_ERROR)\2$(C_RESET)/' \
-e 'p' \
-e 'd' \
'$1' | $(SED) -e 's/\\\\/\\\\\\\\/g'
endef

# Display information about what is being done
# $(call echo-build,<input file>,<output file>,[<run number>])
echo-build	= echo "$(C_BUILD)= $1 --> $2$(if $3, ($3),) =$(C_RESET)"

# Useful shell definitions
sh_true		:= :
sh_false	:= ! :

define get-inputs
$(SED) \
-e '/^INPUT/!d' \
-e 's!^INPUT \(\./\)\{0,1\}!!' \
-e 's/[[:space:]]/\\ /g' \
-e 's/\(.*\)\.aux$$/\1.tex/' \
-e '/\.tex$$/b addtargets' \
-e '/\.cls$$/b addtargets' \
-e '/\.sty$$/b addtargets' \
-e '/\.pstex_t$$/b addtargets' \
-e '/\.dot_t$$/b addtargets' \
-e 'd' \
-e ':addtargets' \
-e 's!.*!$2: &!' \
'$1' | sort | uniq
endef

# $(call get-missing-inputs,<log file>,<target files>)
define get-missing-inputs
$(SED) \
-e '$${' \
-e '  /^$$/!{' \
-e '    H' \
-e '    s/.*//' \
-e '  }' \
-e '}' \
-e '/^$$/!{' \
-e '  H' \
-e '  d' \
-e '}' \
-e '/^$$/{' \
-e '  x' \
-e '  s/^\(\n\)\(.*\)/\2\1/' \
-e '}' \
-e '/^::P\(P\{1,\}\)::/{' \
-e '  s//::\1::/' \
-e '  G' \
-e '  h' \
-e '  d' \
-e '}' \
-e '/^::P::/{' \
-e '  s//::0::/' \
-e '  G' \
-e '}' \
-e 'b start' \
-e ':needonemore' \
-e 's/^/::P::/' \
-e 'G' \
-e 'h' \
-e 'd' \
-e ':needtwomore' \
-e 's/^/::PP::/' \
-e 'G' \
-e 'h' \
-e 'd' \
-e ':needthreemore' \
-e 's/^/::PPP::/' \
-e 'G' \
-e 'h' \
-e 'd' \
-e ':start' \
-e '/^! LaTeX Error: File `/{' \
-e '  b needtwomore' \
-e '}' \
-e '/^::0::\(.*\)/{' \
-e '  s//\1/' \
-e '  /Default extension: /!d' \
-e '  s/.*File `\([^'"'"']*\)'"'"' not found.*/\1/' \
-e '  s/[[:cntrl:]]//' \
-e '  /\.tex/!s/$$/.tex/' \
-e '  s/[[:space:]]/\\ /g' \
-e '  h' \
-e '  s/.*/# MISSING input "&" - (presence of comment affects build)/' \
-e '  p' \
-e '  s/.*//' \
-e '  x' \
-e '  s!^.*!$2: &!' \
-e '  p' \
-e '}' \
-e 'd' \
'$1' | sort | uniq
endef

# Escape dots
# $(call escape-fname-regex,str)
escape-fname-regex	= $(subst /,\\/,$(subst .,\\.,$1))

# This tests whether the build target commands should be run at all, from
# viewing the log file.
# $(call test-log-for-need-to-run,<source stem>)
define test-log-for-need-to-run
$(SED) \
-e '/^No file $(call escape-fname-regex,$1)\.aux\./d' \
$(OUTPUT_DIR)/$1.log \
| egrep -q '^(.*Rerun .*|No file $1\.[^.]+\.|No file .+\.tex\.|LaTeX Warning: File.*)$$'
endef

# check whether the _content_ of two pdf files differs
define pdf-equal
bash -c "[ -f '$(OUTPUT_DIR)/$1.pdf' ] && diff \
	<(grep -a -v 'Type/XRef/Index' '$(OUTPUT_DIR)/$1.pdf' 2> /dev/null | grep -a -v /CreationDate) \
	<(grep -a -v 'Type/XRef/Index' '$(OUTPUT_DIR)/$1.pdf.prev' 2> /dev/null | grep -a -v /CreationDate) \
	> /dev/null 2>&1"
endef

# $(call test-run-again,<source stem>)
test-run-again	= egrep -q '^(.*Warning.*Rerun .*|No file $1.*\.)$$' '$(OUTPUT_DIR)/$1.log'

# Test that a file exists
# $(call test-exists,file)
test-exists		= [ -e '$1' ]

# Note that $(DIFF) returns success when the files are the SAME....
# $(call test-different,sfile,dfile)
test-different		= ! diff -q '$1' '$2' >/dev/null 2>&1
test-exists-and-different	= \
	$(call test-exists,$2) && $(call test-different,$1,$2)

# $(call move-files,source,destination)
move-if-exists		= $(call test-exists,$1) && mv '$1' '$2'

# Copy file1 to file2 only if file2 doesn't exist or they are different
# $(call copy-if-different,sfile,dfile)
copy-if-different	= $(call test-different,$1,$2) && cp '$1' '$2'
copy-if-exists		= $(call test-exists,$1) && cp '$1' '$2'
move-if-different	= $(call test-different,$1,$2) && mv '$1' '$2'
replace-if-different-and-remove	= \
	$(call test-different,$1,$2) \
	&& $(MV) '$1' '$2' \
	|| $(call remove-files,'$1')

# This removes files without checking whether they are there or not.  This
# sometimes has to be used when the file is created by a series of shell
# commands, but there ends up being a race condition: make doesn't know about
# the file generation as quickly as the system does, so $(wildcard ...) doesn't
# work right.  Blech.
# $(call remove-temporary-files,filenames)
remove-temporary-files	= $(if $(KEEP_TEMP),:,$(if $1,rm $1,:))

# Copy a file and log what's going on
# $(call copy-with-logging,<source>,<target>)
define copy-with-logging
if [ -d '$2/' ]; then \
	if cp '$1' '$2/'; then \
		echo "$(C_INFO)Copied '$1' to '$2/'$(C_RESET)"; \
	else \
		echo "$(C_ERROR)Failed to copy '$1' to '$2/'$(C_RESET)"; \
	fi; \
fi
endef

# Makes a an aux file that only has stuff relevant to the target in it
# $(call make-auxtarget-file,<flattened-aux>,<new-aux>)
define make-auxtarget-file
$(SED) \
-e '/^\\newlabel/!d' \
$1 > $2
endef

# Get all important .aux files from the top-level .aux file and merges them all
# into a single file, which it outputs to stdout.
#
# $(call flatten-aux,<toplevel aux>,<output file>)
define flatten-aux
$(SED) \
-e '/\\@input{\(.*\)}/{' \
-e     's//\1/' \
-e     's![.:]!\\&!g' \
-e     'h' \
-e     's!.*!\\:\\\\@input{&}:{!' \
-e     'p' \
-e     'x' \
-e     's/\\././g' \
-e     's/.*/r &/p' \
-e     's/.*/d/p' \
-e     's/.*/}/p' \
-e     'd' \
-e '}' \
-e 'd' \
'$1' > "$1.$$$$.sed.make"; \
$(SED) -f "$1.$$$$.sed.make" '$1' > "$1.$$$$.make"; \
$(SED) \
-e '/^\\relax/d' \
-e '/^\\bibcite/d' \
-e 's/^\(\\newlabel{[^}]\{1,\}}\).*/\1/' \
"$1.$$$$.make" | sort > '$2'; \
$(call remove-temporary-files,$1.$$$$.make $1.$$$$.sed.make)
endef

# Checks for the use of import.sty and bails - we don't support subdirectories
#
# $(call die-on-import-sty,<log file>)
define die-on-import-sty
if egrep -s '/import.sty\)' '$1'; then \
	echo "$(C_ERROR)import.sty is not supported - included files must"; \
	echo "$(C_ERROR)be in the same directory as the primary document$(C_RESET)"; \
	exit 1; \
fi
endef

# Checks for the existence of a .aux file, and dies with an error message if it
# isn't there.  Note that we pass the file stem in, not the full filename,
# e.g., to check for foo.aux, we call it thus: $(call die-on-no-aux,foo)
#
# $(call die-on-no-aux,<aux stem>)
define die-on-no-aux
if [ ! -e '$1.aux' ]; then \
	$(call colorize-latex-errors,$1.log); \
	echo "$(C_ERROR)Error: failed to create $1.aux$(C_RESET)"; \
	exit 1; \
fi
endef


.PHONY: get-temporary-texmf
get-temporary-texmf:
	@$(eval TEXMFHOME ::= $(tmp_texmfhome))

%.$(build_target_extension).1st.make %.aux: %.tex
	@echo "$(C_INFO)Makefile version $(MAKEFILE_VERSION)$(C_RESET)"
	@$(call echo-build,$(OUTPUT_DIR)/$<,$(OUTPUT_DIR)/$*.$(build_target_extension).1st.make,1)
	@$(call run-latex,$<,-recorder $(LATEX_ARGS)) || $(sh_true)
	@cp '$(OUTPUT_DIR)/$*.log' '$(OUTPUT_DIR)/$*.1.log'; \
	$(call die-on-import-sty,$(OUTPUT_DIR)/$*.log); \
	$(call die-on-no-aux,$(OUTPUT_DIR)/$*); \
	$(call flatten-aux,$(OUTPUT_DIR)/$*.aux,$(OUTPUT_DIR)/$*.aux.make)
	@$(call get-inputs,$(OUTPUT_DIR)/$*.fls,$(addprefix $*.,aux aux.make d $(build_target_extension))) >> $(OUTPUT_DIR)/$*.d
	@$(call get-missing-inputs,$(OUTPUT_DIR)/$*.log,$(addprefix $*.,aux aux.make d $(build_target_extension))) >> $(OUTPUT_DIR)/$*.d
	@#  $(SED) -E -i -e '/$(IGNORE_TARGETS)/d' $(OUTPUT_DIR)/$*.d
	@#  egrep -q "# MISSING" $(OUTPUT_DIR)/$*.d && sleep 1 && rm $(OUTPUT_DIR)/$*.pdf ; \
	$(call move-if-exists,$(OUTPUT_DIR)/$*.$(build_target_extension),$(OUTPUT_DIR)/$*.$(build_target_extension).1st.make) || $(sh_true)

# main target
%.$(build_target_extension): %.aux %.$(build_target_extension).1st.make
	@\
	trap "trap - EXIT INT TERM HUP; [ -n '$(tmp_texmfhome)' ] && rm -rf -- '$(tmp_texmfhome)'" EXIT INT TERM HUP; \
	fatal=`$(call colorize-latex-errors,$(OUTPUT_DIR)/$*.log)`; \
	if [ x"$$fatal" != x"" ]; then \
		echo "$$fatal"; \
		exit 1; \
	fi ; \
	run=1; \
	if [ x"$$run" = x"1" -a $(N) != 1 ]; then \
		for i in `seq 2 $(N)`; do \
			$(call echo-build,$*.tex,$(OUTPUT_DIR)/$@,$$i); \
			[ -f '$(OUTPUT_DIR)/$*.pdf' ] && cp '$(OUTPUT_DIR)/$*.pdf' '$(OUTPUT_DIR)/$*.pdf.prev'; \
			$(call run-latex,$*); \
			cp '$(OUTPUT_DIR)/$*.pdf' '$(OUTPUT_DIR)/$*.pdf.prev'; \
			cp '$(OUTPUT_DIR)/$*.log' '$(OUTPUT_DIR)/$*.'$$i'.log'; \
			fatal=`$(call colorize-latex-errors,$(OUTPUT_DIR)/$*.log)`; \
			if [ x"$$fatal" != x"" ]; then \
				echo "$$fatal"; \
				exit 1; \
			fi ; \
			$(call pdf-equal,$*) && ! $(call test-run-again,$*) && break; \
		done; \
	else \
		mv '$(OUTPUT_DIR)/$@.1st.make' '$(OUTPUT_DIR)/$@'; \
	fi; \
	if ! $(call test-exists,$(OUTPUT_DIR)/$*.log) ; then \
		echo "$(C_ERROR)No log file found. This indicates that there is something wrong with your latex installation.$(C_RESET)" ; \
		exit 100 ; \
	fi; \
	$(call copy-with-logging,$(OUTPUT_DIR)/$@,$(BINARY_TARGET_DIR)); \
	$(call latex-color-log,$(OUTPUT_DIR)/$*);


# log target
%.log: %.$(build_target_extension)
	@$(call copy-with-logging,$(OUTPUT_DIR)/$@,$(BINARY_TARGET_DIR))


.PHONY: Targets.ini $(HOME)/.latex-makefile/Targets.ini
-include Targets.ini
-include $(HOME)/.latex-makefile/Targets.ini
#
# vim: noet sts=0 sw=8 ts=8