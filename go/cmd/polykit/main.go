package main

import (
	"fmt"
	"log"
	"os"
	"polysdk/consts"
	"polysdk/internal/config"
	"polysdk/internal/crypto/aesx"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "polykit",
		Usage:       "poly sdk tool kit",
		Description: "polykit is an application to make polysdk config file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:       "filepath",
				Aliases:    []string{"f"},
				Usage:      "config file path",
				Value:      "./polysdk_config.json",
				EnvVars:    []string{consts.EnvPolyConfigPath},
				Required:   true,
				HasBeenSet: true,
			},
			&cli.StringFlag{
				Name:     "description",
				Aliases:  []string{"desc", "d"},
				Usage:    "description of this config file",
				Value:    "",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "initialize a config file",
				Action: func(c *cli.Context) error {
					cfgFile := c.String("filepath")
					cfg := config.NewInitConfig(c.String("description"))
					if err := cfg.StoreFile(cfgFile, false); err != nil {
						return err
					}
					fmt.Println(cfgFile, "ok")
					return nil
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"c"},
				Usage:   "crypt a config file",
				Action: func(c *cli.Context) error {
					cfgFile := c.String("filepath")
					cfg := &config.PolyConfig{}
					if err := cfg.LoadFile(cfgFile); err != nil {
						return err
					}
					if cfg.Key.SecretKey == "" {
						return fmt.Errorf("missing key.secretKey")
					}

					if err := cfg.Validate(); err == nil {
						return fmt.Errorf("it was already crypted")
					}

					keys, err := cfg.GetCryptoKeys()
					if err != nil {
						return err
					}
					crypted, err := aesx.EncodeString(cfg.Key.SecretKey, keys...)
					if err != nil {
						return fmt.Errorf("encrypt fail: %s", err.Error())
					}
					cfg.Key.SecretKey = crypted

					if err := cfg.StoreFile(cfgFile, true); err != nil {
						return err
					}

					fmt.Println(cfgFile, "ok")
					return nil
				},
			},
			{
				Name:    "verify",
				Aliases: []string{"v"},
				Usage:   "verify if a config file is ready",
				Action: func(c *cli.Context) error {
					cfgFile := c.String("filepath")
					cfg := &config.PolyConfig{}
					if err := cfg.LoadFile(cfgFile); err != nil {
						return err
					}
					if err := cfg.Validate(); err != nil {
						return fmt.Errorf("%s %s", cfgFile, err.Error())
					}

					fmt.Println(cfgFile, "ok")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
