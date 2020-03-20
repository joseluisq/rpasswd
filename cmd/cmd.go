package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli/v2"
)

const (
	// DIGITS contains default quantity of digits available
	DIGITS = 10
	// SYMBOLS contains default quantity of symbols available
	SYMBOLS = 30
	// LENGTH contains default characters length
	LENGTH = 40
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	app := &cli.App{
		Name:  "rPasswd",
		Usage: "A secure random password generator tool",
		Description: `CLI tool to generate passwords as described by AgileBits 1Password[1]. The algorithm is commonly used when generating website passwords.

		[1] https://discussions.agilebits.com/discussion/23842/how-random-are-the-generated-passwords`,
		ArgsUsage: `rpasswd [global options] <length>

	 Example: rpasswd --symbols 20 --uppercase 32`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "length",
				Aliases: []string{"l"},
				Value:   LENGTH,
				Usage:   "Password length",
			},
			&cli.IntFlag{
				Name:    "digits",
				Aliases: []string{"d"},
				Value:   0,
				Usage:   fmt.Sprintf("quantity of digits (max. %d)", DIGITS),
			},
			&cli.IntFlag{
				Name:    "symbols",
				Aliases: []string{"s"},
				Value:   0,
				Usage:   fmt.Sprintf("quantity of symbols (max. %d)", SYMBOLS),
			},
			&cli.BoolFlag{
				Name:    "uppercase",
				Aliases: []string{"u"},
				Value:   true,
				Usage:   "allow upper and lower cases",
			},
			&cli.BoolFlag{
				Name:    "repeat",
				Aliases: []string{"r"},
				Value:   false,
				Usage:   "allow characters to repeat",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Value:   false,
				Usage:   "disallow verbose mode",
			},
		},

		Commands: []*cli.Command{
			{
				Name:      "gen",
				Usage:     "Generate a random password including lower-upper cases, digits and symbols characters by default.",
				ArgsUsage: "<length>",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "length",
						Aliases: []string{"l"},
						Value:   LENGTH,
						Usage:   "password length",
					},
					&cli.BoolFlag{
						Name:    "digits",
						Aliases: []string{"d"},
						Value:   true,
						Usage:   fmt.Sprintf("enable %d digit characters", DIGITS),
					},
					&cli.BoolFlag{
						Name:    "symbols",
						Aliases: []string{"s"},
						Value:   true,
						Usage:   fmt.Sprintf("enable %d symbol characters", SYMBOLS),
					},
					&cli.BoolFlag{
						Name:    "uppercase",
						Aliases: []string{"u"},
						Value:   true,
						Usage:   "allow upper-lower case characters. false for lowercase only",
					},
					&cli.BoolFlag{
						Name:    "repeat",
						Aliases: []string{"r"},
						Value:   false,
						Usage:   "allow characters to repeat",
					},
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q"},
						Value:   false,
						Usage:   "disallow verbose mode",
					},
				},
				Action: onGenerateSubCommand,
			},
		},
		Action: onGlobalCommand,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func onGlobalCommand(c *cli.Context) (err error) {
	// 1. process first argument
	len := c.Int("length")

	if c.NArg() > 0 {
		s := c.Args().Get(0)
		len, err = strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	// 2. process flags
	digits := c.Int("digits")
	symbols := c.Int("symbols")
	uppercase := c.Bool("uppercase")
	repeat := c.Bool("repeat")
	quiet := c.Bool("quiet")

	if !quiet {
		fmt.Printf("Generating password %d characters long...\n", len)
	}

	str, err := password.Generate(len, digits, symbols, !uppercase, repeat)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if quiet {
		fmt.Print(str)
	} else {
		fmt.Println("Password generated:")
		fmt.Println(str)
	}

	return nil
}

func onGenerateSubCommand(c *cli.Context) (err error) {
	// 1. process first argument
	len := c.Int("length")

	if c.NArg() > 0 {
		s := c.Args().Get(0)
		len, err = strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	// 2. process flags
	digits := c.Bool("digits")
	symbols := c.Bool("symbols")
	uppercase := c.Bool("uppercase")
	repeat := c.Bool("repeat")
	quiet := c.Bool("quiet")

	if !quiet {
		fmt.Printf("Generating password %d characters long...\n", len)
	}

	digitsN := 0

	if digits {
		digitsN = DIGITS
	}

	symbolsN := 0
	if symbols {
		symbolsN = SYMBOLS
	}

	str, err := password.Generate(len, digitsN, symbolsN, !uppercase, repeat)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if quiet {
		fmt.Print(str)
	} else {
		fmt.Println("Password generated:")
		fmt.Println(str)
	}

	return nil
}
