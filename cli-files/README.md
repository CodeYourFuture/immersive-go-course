# CLI & Files

Become familiar with go and the toolchain by building a JSON-file reading CLI tool: `get-data ".flowers[0].name" garden.json`

Timebox: 4 days

Objectives:

- Install and use [cobra][cobra]
- Use go build/go install/go get etc
- Understand what a process is & the basics of process - lifecycle
- Accept arguments on the CLI
- Open, read (and close) files from CLI arguments
- Learn about JSON and parsing
- Read strings from arguments
- Extract data from a JSON file and print it to the CLI
- Add support for YAML

## Instructions

You're going to build a command-line application that reads data from files.

Let's say you have a `garden.json` file like this, containing data about what's growing in our garden:

```json
{
  "flowers": [
    {
      "name": "Great Maiden's Blush",
      "genus": "Rosa",
      "species": "Rosa × alba",
      "color": "white",
      "height": 2.2
    },
    {
      "name": "Bright Gem",
      "genus": "Tulipa",
      "species": "Tulipa linifolia",
      "color": "Yellow/orange",
      "height": 0.17
    },
    {
      "name": "Elizabeth",
      "genus": "Magnolia",
      "species": "M. acuminata × M. denudata",
      "color": "Pale yellow",
      "height": 9.6
    }
  ]
}
```

The tool you will build will extra data from this file. Here are some examples:

```bash
> get-data ".flowers[0].name" garden.json
"Great Maiden's Blush"

> get-data ".flowers[1].height" garden.json
2.2

> get-data ".flowers[2]" garden.json
{"name": "Elizabeth", "genus": "Magnolia", "species": "M. acuminata × M. denudata", "color": "Pale yellow", "height": 9.6 }
```

[cobra]: https://github.com/spf13/cobra#overview
