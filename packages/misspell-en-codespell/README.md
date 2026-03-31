---
package_name: misspell-en-codespell
package_title: Misspell EN Codespell
package_desc: Replace commonly misspelled english words (codespell).
package_version: 0.1.0
package_author: Timo Runge
package_repo: https://github.com/timorunge/espanso-misspell
---
# misspell-en-codespell

misspell-en-codespell is an espanso package which is replacing commonly
misspelled english words. The package is auto-generated from the
[codespell](https://github.com/codespell-project/codespell) dictionary,
deduplicated against misspell-en to avoid trigger conflicts.

## Installation

```
espanso install misspell-en-codespell
espanso restart
```

## Usage

Type `acessibility` and see it replaced with `accessibility`.

## License

[CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/)

Dictionary data from codespell is derived from English Wikipedia and licensed
under [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/).