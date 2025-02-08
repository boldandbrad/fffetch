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

> Coming soon.

## Dev Setup

Clone this repo:

```
git clone https://github.com/boldandbrad/fffetch
```

Ensure [Go](https://go.dev/) is installed, matching or exceeding the version
listed in [go.mod](./go.mod).

Then, run the project:

```
go run .
```

By default it will attempt to fetch and parse data for all NFL teams from the
most recent season. In order to avoid rate limits on Pro Football Reference,
requests are spaced out by roughly 5 seconds each, which means the script may
run for a couple minutes. Canceling the job at any time is OK.

Modify the default behavior by updating the `teams` and `years` slices in
[fffetch.go](./fffetch.go).

## License

[MIT](./LICENSE)
