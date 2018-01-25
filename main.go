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
	app.Version = "0.1"
	app.Name = "barcoder"
	app.Usage = "Generate a label (usually with a barcode)"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "barcoder.toml",
			Usage: "Load configuration from FILE",
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
		return barcodes.Process(*config, output, c.Args())
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

}

func loadConfig(filename string) (*barcodes.BarcodeConfig, error) {
	if content, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		var config barcodes.BarcodeConfig
		if err := toml.Unmarshal(content, &config); err != nil {
			return nil, err
		}
		return &config, nil
	}

}
