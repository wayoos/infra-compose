# infra-compose

## Setting up Go to work on infra-compose

If you have never worked with Go before, you will have to complete the following
steps in order to be able to compile and test infra-compose. These instructions target MacOS environments
so you may need to adjust them for Linux, Windows or other shells.

1. install [go](https://golang.org) and [dep](https://golang.github.io/dep/), with [Homebrew](https://brew.sh).

    ```
    brew install go
    brew install dep
    ```

2. Set and export the `GOPATH` environment variable and update your `PATH`. For
   example, you can add to your `.bash_profile`.

    ```
    export GOPATH=$HOME/<path to your go workspace>
    export GOROOT=/usr/local/opt/go/libexec
    export PATH="$PATH:${GOPATH}/bin:${GOROOT}/bin"
    ```
    
    Find GOROOT with `brew --prefix go` and add 'libexec'

3. Download the infra-compose source (and its dependencies) by running `go get
   github.com/wayoos/infra-compose`. This will download the infra-compose source to
   `$GOPATH/src/github.com/wayoos/infra-compose`.
