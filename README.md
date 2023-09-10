# paste

Yet anothe paste bin.

## Features

- written in Go
- self-contained binary (except the database, if used), no need for `cgo`
- can be used with SQLite, including in-memory setup (use `DB_DSN=:memory:`)
- can (probably) be used witn PostgreSQL, which is supported, but wasn't tested (please open a PR if you find a bug)
- uses [htmx](https://htmx.org/) and [TailwindCSS](https://tailwindcss.com/)
- syntax highlighting using [chroma](https://github.com/alecthomas/chroma) (the same one used in Hugo)
- resulting pastes can be suffixed with the extension (`.md`, `.py`, etc.) and highlighted accordingly
- the contents are selectable with `Ctrl+A` and line-numbers will not be copied (see tips & tricks if you are using Firefox)

## Self-hosting

The paste server can be configured with these environment variables:

```bash
BASE_URL=http://my.paste.com
DB_DSN=file:paste.db
BIND=:3000
```

And here is an example of running paste server in a docker container:

```
docker run -p 3002:3000 \
--env BASE_URL=http://localhost:3000 \
--env DB_DSN=file:paste.db \
--env BIND=:3000 \
--volume $PWD/tmp.db:/app/paste.db \
paste:latest server
```

If you want to map the file from the container to a local filesystem, make sure to create it beforehand `touch paste.db`. (please open a PR if you would like to fix that)

## Development

- `task deps` - `tidy` and `vendor`
- `task build-tailwind` - rebuild css files, and put them in `web/static/output.css`

> **Note**: this repository contains all files required for running the server.
> This includes both `vendor` and built frontend files.

- `task dev --watch` - the easiest way to run backend and frontend with auto reloading on file changes

## Tips and tricks

### Command-line usage

If you want to use this utility from the command line, then the fastest way is to make an alias:

```bash
export PASTE_URL=http://<my-paste-instance>
alias pb="curl --data-urlencode content@- $PASTE_URL"
```

You then use it like so: `cat README.md | pb`. It will output the url of the paste.

The ultimate workflow on a Mac is to copy the url right after pasting, so the command would look like `cat README.md | pb | pbcopy`.

### Dealing with Firefox copy-paste extra lines bug

There is [bug](https://bugzilla.mozilla.org/show_bug.cgi?id=1273836) in Firefox from 2016 that adds extra white space on copied text when the page contains `user-select: none` styles.

If you add `.txt` to the paste - http://<my-paste-url>/wrzivrbrki.txt - then it will output plain text, which can be copied from Firefox without extra new lines.

This issue doesn't exist in Chrome and Chromium-based browsers.

## Honorable mentions

- [w4/bin](https://github.com/w4/bin) - for inspiration

## Contributing

You are welcome to open new issues and submit PRs.

Here are some ideas:

- move more settings into env vars (default name length and maximum file size)
- test with PostgreSQL and let me know how it went in a new issue
- fix an issue when a directory is created instead of a file when using SQLite file and a volume

