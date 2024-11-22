# Contributing to urlshortener

Thank you for considering contributing to the project! We welcome contributions in all forms.

## Code of Conduct

Please adhere to the [Go Community Code of Conduct](https://go.dev/conduct) when interacting with others in the project.

## How to Contribute

1. Fork the repository.
2. Create a new branch (`git checkout -b my-feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin my-feature-branch`).
6. Create a new Pull Request.

## Build

### Docker

To build the Docker image, use the following command:

````bash
docker build .
````

## Testing

### Test code

```sh
make test
# will run `go test ./... -parallel 10`
```

## Issue Reporting

If you encounter a bug or have a feature request, please open an issue in the GitHub repository.
