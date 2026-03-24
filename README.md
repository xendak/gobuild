# gobuild

Will eventually be a emacs like compilation-mode, but in tui, so its editor
agnostic.Done in go, so that i can both learn the language, and make use of
bubbletea, which sounds a really cool library name

### Usage

```bash
go build . 2>&1 | gobuild
gobuild 'somefile'
```
