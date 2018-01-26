Barcoder
========

Generate ZPL for Zebra label printers. Generate Barcodes, QRCodes and import
images.

Under active delopment.

Installation
------------

    go get github.com/codenaut/barcoder


Usage
-----

    NAME:
    barcoder - Generate a label (usually with a barcode) and generate output in ZPL for Zebra label printers

    USAGE:
        barcoder [global options] command [command options] [arguments...]

    VERSION:
        0.1

    COMMANDS:
        help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
        --config value, -c value  Load configuration from FILE (default: "barcoder.toml")
        --output value, -o value  Output to FILE. Default to std output
        --help, -h                show help
        --version, -v             print the version


Previewing ZPL files
--------------------
Labelary has brilliant online previewer: http://labelary.com/viewer.html
