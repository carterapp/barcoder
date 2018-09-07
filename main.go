package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
	"gopkg.in/urfave/cli.v1"

	"github.com/codenaut/barcoder/barcodes"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.2"
	app.Name = "barcoder"
	app.Usage = "Generate a label (usually with a barcode) and generate output in ZPL for Zebra label printers"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "x-offset, x",
			Value: 0,
			Usage: "X Offset",
		},
		cli.IntFlag{
			Name:  "y-offset, y",
			Value: 0,
			Usage: "Y Offset",
		},

		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "Load configuration from FILE. Default to std input",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output to FILE. Default to std output",
		},
	}
	app.Action = func(c *cli.Context) error {
		config, err := loadConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}
		output := os.Stdout
		outputfile := c.GlobalString("output")
		if outputfile != "" {
			output, err = os.Create(outputfile)
			if err != nil {
				return err
			}
		}
		xOffset := c.GlobalInt("x-offset")
		yOffset := c.GlobalInt("y-offset")

		processor := barcodes.New(*config, xOffset, yOffset)
		return processor.Process(output, c.Args())
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

}

func loadConfig(filename string) (*barcodes.BarcodeConfig, error) {
	unmarshall := func(content []byte) (*barcodes.BarcodeConfig, error) {
		var config barcodes.BarcodeConfig
		if err := toml.Unmarshal(content, &config); err != nil {
			return nil, err
		}
		return &config, nil
	}
	if filename == "" {
		if content, err := ioutil.ReadAll(os.Stdin); err != nil {
			return nil, err
		} else {
			return unmarshall(content)
		}
	} else if content, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return unmarshall(content)
	}

}
