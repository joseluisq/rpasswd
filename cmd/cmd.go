package cmd

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	"github.com/miquella/ask"
)

const (
	// DIGITS contains default quantity of digits available
	DIGITS = 10
	// SYMBOLS contains default quantity of symbols available
	SYMBOLS = 30
	// LENGTH contains default characters length
	LENGTH = 40
)

const (
	algoBcrypt = "bcrypt"
	algoScrypt = "scrypt"
	algoArgon2 = "argon2"
	algoPbkdf2 = "pbkdf2"
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
			{
				Name:      "enc",
				Usage:     "Encrypt a given password using a key derivation function hash.",
				ArgsUsage: "<hash function>",
				Action:    onEncryptSubCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hash",
						Aliases: []string{"s"},
						Value:   "pbkdf2",
						Usage:   "key derivation function hash like bcrypt, scrypt, argon2 or pbkdf2",
					},
				},
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

func onEncryptSubCommand(c *cli.Context) (err error) {
	hash := c.String("hash")

	if c.NArg() > 0 {
		hash = c.Args().Get(0)
	}

	supportedHash := false

	switch hash {
	case algoBcrypt, algoScrypt, algoArgon2, algoPbkdf2:
		supportedHash = true
	}

	if !supportedHash {
		return fmt.Errorf("`%s` is not a supported key derivation function hash", hash)
	}

	stdinPass, err := ask.HiddenAsk("New secure password: ")

	if err != nil {
		return err
	}

	stdinSalt, err := ask.HiddenAsk("A random salt (at least 8 bytes): ")

	if err != nil {
		return err
	}

	passwdEnc := hashPassword(stdinPass, stdinSalt, hash)

	fmt.Println("Password encrypted:")
	fmt.Println(passwdEnc)

	return nil
}

func hashPassword(passwd, salt, algo string) string {
	var tmpPasswd []byte

	switch algo {
	case algoBcrypt:
		tmpPasswd, _ = bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
		return string(tmpPasswd)
	case algoScrypt:
		tmpPasswd, _ = scrypt.Key([]byte(passwd), []byte(salt), 65536, 16, 2, 50)
	case algoArgon2:
		tmpPasswd = argon2.IDKey([]byte(passwd), []byte(salt), 2, 65536, 8, 50)
	case algoPbkdf2:
		fallthrough
	default:
		tmpPasswd = pbkdf2.Key([]byte(passwd), []byte(salt), 10000, 50, sha256.New)
	}

	return fmt.Sprintf("%x", tmpPasswd)
}
