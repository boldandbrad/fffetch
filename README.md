# fffetch

> **F**antasy **F**ootball data **F**etcher

Fetch and analyze NFL player data from
[Pro Football Reference](https://www.pro-football-reference.com/) for use in
personal Fantasy Football simulations.

> [!CAUTION]
> Use at your own risk. This project utilizes web-scraping techniques to
> retrieve data from Pro Football Reference, which may be against their terms of
> service. The maintainers of fffetch are not responsible for your use of this
> tool.

## Usage

Build the binary:

```bash
go build -o fffetch
```

Run the fetch command:

```bash
./fffetch fetch
```

By default, this fetches data for all NFL teams from the previous season.

### Options

- `-t, --team <team>`: Specify teams to fetch (e.g., `-t KC`, `-t BUF -t PHI`). Defaults to all teams.
- `-y, --year <year>`: Specify years to fetch (e.g., `-y 2023`, `-y 2023 -y 2024`). Defaults to previous year.
- `-f, --force`: Force re-fetch existing data instead of skipping it.

### Examples

Fetch Lions 2023 data:

```bash
./fffetch fetch -t DET -y 2023
```

Fetch multiple teams and years:

```bash
./fffetch fetch -t DET -t GB -t MIN -y 2023 -y 2024
```

Re-fetch existing data:

```bash
./fffetch fetch --force
```

### Output

The tool displays an interactive progress bar (in supported terminals) with status updates for each team/year combination. Data is saved to CSV files in the `output/final/` directory.

### Notes

- Requests are spaced out by 2-4.5 seconds to avoid rate limiting on Pro Football Reference
- Canceling the job at any time is OK (press `q` or `Ctrl+C` in interactive mode)
- Existing data is automatically skipped unless using `--force` flag

## Dev Setup

Clone this repo:

```
git clone https://github.com/boldandbrad/fffetch
```

Ensure [Go](https://go.dev/) is installed, matching or exceeding the version
listed in [go.mod](./go.mod).

Then, build and run:

```bash
go build -o fffetch
./fffetch fetch
```

## License

[MIT](./LICENSE)
