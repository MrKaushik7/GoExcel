# GoExcel

* A simple implementation of Excel in GoLang without the UI
* batch program which accepts a CSV input and a stream of commands, that outputs a modified CSV file.
* Example ->
```csv
A | B
1 | 2
3 | 4
=A1+B1 | =A2+B2
```

```csv
A | B
1 | 2
3 | 4
3 | 7
```

* Parses -
- Mathematical Expressions
- Formulae
* Detects and reports formulaic cycles

Internally separated by '|' to avoid function argument conflicts.