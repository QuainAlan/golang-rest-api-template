# Contributing

Thank you for helping improve this template. Contributions are welcome via issues and pull requests.

## Before you start

- Check existing [issues](https://github.com/LAA-Software-Engineering/golang-rest-api-template/issues) and pull requests to avoid duplicate work.
- For larger changes, open an issue first so maintainers can agree on direction (see also the [README](./README.md)).

## Development setup

- **Go**: version declared in [`go.mod`](./go.mod) (module uses vendoring; prefer builds with `-mod=vendor` if that is your local default).
- **Docker** (optional): for full stack parity with [`docker-compose.yml`](./docker-compose.yml).

Useful commands:

```bash
go test ./... -race
```

```bash
make test          # tests with coverage flags per Makefile
make up            # run stack via Docker Compose
```

## Pull requests

- Keep changes focused on one concern when possible.
- Run **`go test ./... -race`** before opening a PR.
- Match existing style (`gofmt`, package layout, naming).
- If you change HTTP routes, request/response shapes, or models used in Swagger, regenerate API docs (see the `setup` target in the [`Makefile`](./Makefile) and `swag init` usage there) and include updated `docs/` artifacts if applicable.
- Do not commit secrets, real credentials, or personal `.env` files.
- New or changed environment variables should be reflected in the README (and `.env.example` when that file exists).

Pull requests use the template under [`.github/pull_request_template.md`](.github/pull_request_template.md) when opened on GitHub.

## Conduct

This project follows the [Code of Conduct](./CODE_OF_CONDUCT.md). By participating, you agree to uphold it.

## License

By contributing, you agree that your contributions are licensed under the same terms as the project ([MIT License](./LICENSE)).
