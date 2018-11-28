# infra-compose

## Installing

There are different ways to get going install infra-compose:

Using homebrew

```
brew install wayoos/tap/infra-compose
```

Manually

Download your preferred flavor from the [releases page](https://github.com/wayoos/infra-compose/releases/latest) and install manually.

## Setting up Go to develop on infra-compose

If you have never worked with Go before, you will have to complete the following
steps in order to be able to compile and test infra-compose. These instructions target MacOS environments
so you may need to adjust them for Linux, Windows or other shells.

1. install [go](https://golang.org), with [Homebrew](https://brew.sh).

    ```
    brew install go
    ```

2. Download the infra-compose source by running `git clone https://github.com/wayoos/infra-compose.git`.
   This will download the infra-compose source.

3. Launch the application with `go run main.go`

## Release

Release binaries are deployed to Github.

1. GoReleaser will use the latest Git tag of your repository. Create a tag and push it to GitHub:

```
$ git tag -a v0.1.0 -m "First release"
$ git push origin v0.1.0
```

2. [Create a github token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)
    and set it as the env variable `GITHUB_TOKEN`. `github-release` will automatically pick it up from the
    environment so that you don't have to pass it as an argument.

3. Execute the command `./build.sh release`

## Variable

https://github.com/gliderlabs/sigil