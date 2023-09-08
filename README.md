# paste

Yet anothe paste bin alternative.

## Development

- `task deps` - `tidy` and `vendor`
- `task build-tailwind` - rebuild css files, and put them in `web/static/output.css`

> **Note**: this repository contains all files required for running the server.
> This includes both `vendor` and built frontend files.

- `task dev --watch` - the easiest way to run backend and frontend with auto reloading on file changes

## FAQ

- Why plain text for view without extension?
  - There is [bug](https://bugzilla.mozilla.org/show_bug.cgi?id=1273836) in Firefox from 2016 that adds extra white space on copied text when used with `user-select: none`.
    And because my default browser is Firefox, I prefer to be able to copy text from the default view.
    You can still add extension, even `txt`, to show line numbers and nicer formatting.
