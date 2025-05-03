# Contributing

## First setup

1. Install required tools.
   - Go@^1.23 ([Installation](https://go.dev/doc/install))
   - Deno@^2.0 ([Installation](https://deno.com/))
2. [Fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo)
   and
   [clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository)
   the repository.
3. Open it in your favorite editor. [VSC](https://code.visualstudio.com/) is
   recommended. For small changes you can also use
   [in-browser GitHub editor](https://docs.github.com/en/codespaces/the-githubdev-web-based-editor).
4. Open terminal. You can use built-in
   [VSC terminal](https://code.visualstudio.com/docs/terminal/getting-started).
5. Run `deno task init:build` to create required files.
6. Change the `.env` file.
   - Set up JWT secret.
7. Run `deno task serve` to start the server.

## Changing the code base

Client code base (./client) is not tied with the server code base (./server).
The best way is to use 2 terminals (3-rd for other tasks):

> [!NOTE]
> You can use Visual Studio Code's task commands: `Tasks: Run Task`.
>
> - Compile Client & Watch
> - Serve

```bash
deno task compile:client watch
```

```bash
deno task serve
```

> [!WARNING]
> The `serve` script can ignore new files, so it should be started after
> `compile:client` script generates all files. If you are using `watch`, wait
> for "watching..." message.

## How to write commit messages

We use [Conventional Commit messages](https://www.conventionalcommits.org/) to
automate version management.

Most common commit message prefixes are:

- `fix:` which represents bug fixes and generate a patch release.
- `feat:` which represents a new feature and generate a minor release.
- `chore:` which represents a development environment change and generate a
  patch release.
- `docs:` which represents documentation change and generate a patch release.
- `style:` which represents a code style change and generate a patch release.
- `test:` which represents a test change and generate a patch release.
- `BREAKING CHANGE:` which represents a breaking change and generate a major
  release. Or you are able to use `!` at the end of the prefix. For example
  `feat!: new feature` or `fix!: bug fix`.

## Compilation

Go allows us to compile server code, and we are using this feature very well,
providing even additional feature: client embedding.

Embedded client files are stored inside server's binary optionally. You can use
go build tags for setting things for your needs.

Available go build tags:

- Environment: `test`, `prod`. If not provided, `dev`.
  - `test` exists, but we don't have test at this moment.
  - `dev` enables client files watching.
  - `prod` normal mode.
- Client embedding: `lite`. If not provided, `normal`.
  - `normal` enables client files embedding. The server binary will become
    standalone.
  - `lite` disables files embedding. The server binary will use closest
    ./client/static and ./client/templates directories. This option makes the
    server binary more flexible and reduces its size.

Available deno tasks for server compilation:

```bash
deno task compile:server
# go build -o ./dist/server main.go

deno task compile:server:lite
# go build -tags lite -o ./dist/server lite.go

deno task compile:server:prod
# go build -tags prod -o ./dist/server lite.go

deno task compile:server:test
# go build -tags lite test -o ./dist/server lite.go
```

## About DOM (HTMX, Shoelace) and Session

Resources:

- <https://shoelace.style>
- <https://htmx.org/docs/>
- <https://htmx.org/reference/>
- <https://pkg.go.dev/html/template>
- <https://docs.gofiber.io/next/> - v3 (Next), not v2!

We are using HTMX. JavaScript (TypeScript) is an utility for importing
libraries, extending DOM and web-components functionality. We are fetching HTML
from the server instead of JSON.

The session stored in cookies and should be changed this way:

1. Client sends request to change own cookies.
2. Server responds with new cookies.

### About templates

Files in the [./client/templates](./client/templates) can be rendered through
Go's template language: <https://pkg.go.dev/html/template>.

That means, you can use specific syntax and replacements, but the variables
should be declared by the server. You can find more it in the
[./server](./server). Specific functions are declared in the
[./server/engine.go](./server/engine.go). .
