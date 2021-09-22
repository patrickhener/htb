# htb

Hack The Box Report Writer Template und Scripte

## Installation

`git clone git@git.syss.intern:phener/htb.git`

- $HTBDIR definieren
- optional: $HTBAUTHOR definieren
- optional: $HTBPROFILEID definieren

Alles kann zum Beispiel in der .bashrc erledigt werden. Auszug aus meiner .bashrc:

```sh
export HTBDIR="$HOME/htb"
export HTBAUTHOR="C1sc0"
export HTBPROFILEID="34604"
```

Beispiel: klont ihr das Repo nach /home/user/htb, dann ist das euer $HTBDIR und die Unterordner _loot_, _report_, _app_ und _template_ liegen dann darunter.

## Requirements

- pandoc
- pandoc-latex-color
- texlive
- eisvogel latex template
- phantomjs
- obsidian

Unter Arch bekomme ich _pandoc_ und texlive über das repo und _pandoc-latex-color_ kommt über pip.
Eisvogel template findet man hier: [Eisvogel Latex](https://github.com/Wandmalfarbe/pandoc-latex-template/releases/latest)

## Framework app

Die Framework app im Ordner _app_ ist in *go* geschrieben und kann mittels `make build` und `make install` gebaut und in den GOPATH verschoben werden.

Die App kennt verschiedene Befehle:

`htb <mode> <boxname>`

### Create

Erstellt die Ordnerstruktur und kopiert das angepasste Template in den Report Ordner.

### Edit

Öffnet die Markdown file mittels `xdg-open` zum Editieren in Obsidian. WICHTIG! - Man muss aktuell den Vault einmalig von Hand in Obsidian öffnen, da der Automatismus sonst nicht funktioniert.

### Convert

Erstellt aus der Markdown-Datei eine PDF-Datei.

### Open

Öffnet die erstellte PDF mittels `xdg-open`.

### List

Listet alle bereits erstellten Boxen auf.

### Clear

Löscht nach Rückfrage den loot und den report Ordner der Box.

