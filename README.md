# htb

Hack The Box Framework for using 'writeup' Latex style.

## Installation

`git clone git@git.syss.intern:phener/htb.git`

- define $HTBDIR
- define $HTBAUTHOR
- define $HTBPROFILEID

Can easily be done in a .bashrc for example:

```sh
export HTBDIR="$HOME/htb"
export HTBAUTHOR="c1sc0"
export HTBPROFILEID="34604"
```
## Requirements

- texlive
- writeup latex style [github.com/patrickhener/writeup](https://github.com/patrickhener/writeup)
- phantomjs
- VSCode (edit mode calls 'code %s')

## App

App can be found from folder _app_ and is written in *go*. One can build and install it via `make build` and `make install`.

There are different modes:

`htb <mode> <boxname>`

### Create

This mode creates the folder structure for loot and for the writeup and copies over the template from `writeup` Latex style.

### Edit

This mode opens the writeup path of the box in VSCode.

### Open

This mode opens the boxes PDF via `xdg-open`.

### List

This mode lists all boxes in the writeup directory.

### Clear

This mode deletes the corresponding loot and writeup folder of the box after asking.

### Badge

This mode updates your badge.png and copies it over to the 'writeup/images' folder. This will be run automatically everytime you are creating a box.