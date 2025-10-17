# lga — Log Analyser (CLI)
**(Part of a series of mini-projects to get deeper knowledge of Go and its ecosystem)**

A minimal command-line utility to parse and analyse log files. It identifies errors and extracts insights from raw logs. Built with [Cobra](https://github.com/spf13/cobra).

```bash
git clone https://github.com/danilobml/lga.git
cd lga
go build -o lga ./cmd/lga
./lga --help
```

If you prefer a single command installation, you can run:
```bash
go install github.com/danilobml/lga/lga/cmd/lga@latest
```
Ensure that your `$GOPATH/bin` or `$GOBIN` is part of your system `PATH`.

To run the analyser, execute:
```bash
lga [flags] <logfile>
```
If no arguments are provided, `lga` displays help information.

### Examples

```bash
lga ./logs/app.log
lga --paths ./logs/access.log
lga --from "2025-10-14" --to "2025-10-17" ./logs/app.log
lga --from "2025-10-14 08:00:00" --to "2025-10-14 18:00:00" ./logs/app.log
```

### Flags

```
--paths      Enables per-path analysis.
--from       Defines a starting point for filtering logs (inclusive).
--to         Defines an endpoint for filtering logs (exclusive if only a date is provided).
-t, --toggle Example flag scaffolded by Cobra.
-h, --help   Displays the command help.
```

Each flag modifies the analysis scope; the `<logfile>` argument is required and should be a valid path to your log file.

Example log formats:
```
127.0.0.1 - - [14/Oct/2025:09:12:33 +0000] "GET /api/health HTTP/1.1" 200 123 "-" "curl/8.4.0"
127.0.0.1 - - [14/Oct/2025:09:13:05 +0000] "GET /api/users HTTP/1.1" 500 0 "-" "curl/8.4.0"
```

### Date and Time Handling

`--from` and `--to` are optional. When omitted, all entries are considered. When both are provided, only logs within `[from, to)` are analysed. If `--to` is supplied with a date only (no time), that day is excluded. Example:
```
--to "2025-10-17" excludes 2025-10-17 itself (i.e., up to 2025-10-16 23:59:59)
```

### Development

```
.
├─ cmd/
│  └─ lga/          # main package
├─ internal/
│  └─ analyser/     # parsing & analysis logic
├─ go.mod
└─ README.md
```

Run:
```bash
go run ./cmd/lga --help
go run ./cmd/lga --paths --from "2025-10-14" ./logs/app.log
```

Build:
```bash
go build -o lga ./cmd/lga
```

Install locally:
```bash
go install ./cmd/lga
```

Test:
```bash
go test ./...
