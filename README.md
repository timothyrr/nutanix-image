# nutanix-image -- a tool for uploading or downloading nutanix images

## Table of Contents

* [About](#about)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Usage](#usage)
* [Output](#output)
* [Contributing](#contributing)
* [Contact](#contact)


## About

`nutanix-image` is a simple utility for either downloading or uploading any nutanix image file.

Configuration may be provided via a config file (defaults to `~/.nutanix-image.yaml`) or via commandline switches.

An example configuration file that defines all configuration for all commands:

```
---
endpoint: 'https://nutanix.example.com'
username: 'nutanix'
password: 'nutanixs3cr3t'
insecure: true
```


### Built With
* [Cobra](https://github.com/spf13/cobra)
* [Viper](https://github.com/spf13/viper)
* [Resty](https://github.com/go-resty/resty)

## Getting Started

Simple compilation can be done via `go build`
Task builds can be performed via `task build`
Package as an RPM for installation into `/usr/local/bin` with `task build-rpm`


### Prerequisites

* Go v1.12+ is recommended for Go Module support
* fpm (for easy RPM packaging)

[https://fpm.readthedocs.io/en/latest/installing.html](Installation)

* go-task

[https://taskfile.dev/#/installation](Installation)


### Installation

Ensure Golang 1.11+ is installed and GO111MODULE is set to `on`.

Clone the repo
```sh
$ git clone https://github.com/timothyrr/nutanix-image.git
$ cd ./nutanix-image
```

[Method 1] Manually compile and move into $PATH
```sh
$ go build -o nutanix-image (if building on a macOS system to be used with linux use this command: env GOOS=linux GOARCH=amd64  go build -o nutanix-image)
$ mv ./nutanix-image /usr/local/bin/ # or ~/bin, etc. or just execute directly e.g. ./nutanix-image --help
```

[Method 2] Use go-task to install into `/usr/local/bin`
```sh
$ task install
```

## Usage

The following subcommands are available for deleting various objects:

```
Available Commands:
  get     Download a given image
  create  Upload a given image
  help    Help about any command
```

Calling each subcommand with no arguments or with `--help` will display the appropriate help text for the subcommand.

## Output

By default `nutanix-image` will only output errors and will exit with a status of `0` if the desired objects are download/uploaded.

Informational messages may be enabled by passing `-v` or `--verbose` and full debug details of every RESTful HTTP interaction can be enabled with `-d` or `--debug`.

Example `--verbose` output:

```sh
~]$ nutanix-image get rhel8-template-20191001-NY2 -v
INFO[0000] Using config file: ~/.nutanix-image.yaml
INFO[0000] Locating Nutanix image (rhel8-template-20191001-NY2)...
INFO[0002] Found Nutanix image with UUID: 3a043187-b351-463b-acec-beb1122a30a1
INFO[0002] Downloading Nutanix image file to ~/rhel8-template-20191001-NY2.img...
INFO[0030] Nutanix image downloaded

~]$ nutanix-image create rhel8 --source=/var/images/rhel8.img -v
INFO[0000] Using config file: ~/.nutanix-image.yaml
INFO[0002] Creating Nutanix image (rhel8)...
INFO[0003] Uploading Nutanix image file (rhel8.img)...
INFO[0014] Nutanix image uploaded (UUID: 0b295ba0-2b78-4d71-8de7-793d8a44ed0d)
```


### Help docs
```sh
$ nutanix-image
Nutanix-image enables the user to be able to interact with
the nutanix images API for the following:

  * downloading images
  * uploading images

Usage:
  nutanix-image [command]

Available Commands:
  create      Create and upload a nutanix image
  get         Download a nutanix image
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.nutanix-image.yaml)
  -d, --debug           set debug output
  -h, --help            help for nutanix-image
  -v, --verbose         set verbose output
      --version         version for nutanix-image

Use "nutanix-image [command] --help" for more information about a command.
```

```sh
$ nutanix-image get --help
This command will download a given image by name
using the filepath you specify. If no filepath is specified
it will use a default path of your current working directory.

Example:

  $ nutanix-image get rhel8 --output_dir=/var/images

Usage:
  nutanix-image get <image-name> [flags]

Flags:
      --endpoint string     the endpoint exposing the Nutanix API (Ex. https://nutanix.example.org)
  -h, --help                help for get
      --insecure            enable or disable server certificate validation
      --output_dir string   the directory to save the image to (default "/var/images")
      --password string     the password for Nutanix API authentication
      --username string     the username for Nutanix API authentication

Global Flags:
      --config string   config file (default is $HOME/.nutanix-image.yaml)
  -d, --debug           set debug output
  -v, --verbose         set verbose output
```

```sh
$ nutanix-image create --help
This command will create and upload a given image
by name using the filepath you specify with the --source flag.

Example:

  $ nutanix-image create rhel8 --source=/var/images/rhel8.img

Usage:
  nutanix-image create <image-name> --source=<source-file> [flags]

Flags:
      --endpoint string   the endpoint exposing the Nutanix API (Ex. https://nutanix.example.org)
  -h, --help              help for create
      --insecure          enable or disable server certificate validation
      --password string   the password for Nutanix API authentication
      --source string     the source image to upload
      --username string   the username for Nutanix API authentication

Global Flags:
      --config string   config file (default is $HOME/.nutanix-image.yaml)
  -d, --debug           set debug output
  -v, --verbose         set verbose output
```


### Examples

Download an image:
```sh
$ nutanix-image get rhel8
```

Download an image while specifying output directory:
```sh
$ nutanix-image get rhel8 --output_dir=/var/images
```

Download an image while overriding the config file endpoint (but keeping other nutanix settings from file):
```sh
$ nutanix-image get --endpoint https//other-nutanix.example.com rhel8
```

Upload an image:
```sh
$ nutanix-image create rhel8 --source=/var/images/rhel8.img
```


## Contributing

Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
