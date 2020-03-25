package cmd

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	"github.com/miquella/ask"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	app := &cli.App{
		Name:                 "rPasswd",
		EnableBashCompletion: true,
		Usage:                "A small and secure password generator and encryptor tool",
		Description: `CLI tool to generate secure passwords as described by AgileBits 1Password[1]. With support for password encryption as well using bcrypt, scrypt, argon2 or pbkdf2.

		[1] https://discussions.agilebits.com/discussion/23842/how-random-are-the-generated-passwords`,
		ArgsUsage: `rpasswd [global options] <length>

	 Example: rpasswd --symbols 20 --uppercase 32`,
		Flags: []cli.Flag{
			VersionFlag(),
		},

		Commands: []*cli.Command{
			{
				Name:      "gen",
				Aliases:   []string{"g"},
				Usage:     "Generate a random password including lower-upper cases, digits and symbols characters by default",
				ArgsUsage: "<length>",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "length",
						Aliases: []string{"l"},
						Value:   defaultLength,
						Usage:   "password length",
					},
					&cli.BoolFlag{
						Name:    "digits",
						Aliases: []string{"d"},
						Value:   true,
						Usage:   fmt.Sprintf("enable %d digit characters", defaultDigits),
					},
					&cli.BoolFlag{
						Name:    "symbols",
						Aliases: []string{"s"},
						Value:   true,
						Usage:   fmt.Sprintf("enable %d symbol characters", defaultSymbols),
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
				Name:      "cgen",
				Aliases:   []string{"cg"},
				Usage:     "Generate a custom random password with specified lower-upper cases, digits and symbols characters",
				ArgsUsage: "<length>",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "length",
						Aliases: []string{"l"},
						Value:   defaultLength,
						Usage:   "Password length",
					},
					&cli.IntFlag{
						Name:    "digits",
						Aliases: []string{"d"},
						Value:   0,
						Usage:   fmt.Sprintf("quantity of digits (max. %d)", defaultDigits),
					},
					&cli.IntFlag{
						Name:    "symbols",
						Aliases: []string{"s"},
						Value:   0,
						Usage:   fmt.Sprintf("quantity of symbols (max. %d)", defaultSymbols),
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
				Action: onCustomGenerateSubCommand,
			},
			{
				Name:      "enc",
				Usage:     "Encrypt a given password using a hashing function (bcrypt, scrypt) or a key derivation function hash (argon2, pbkdf2)",
				ArgsUsage: "<algorithm>",
				Action:    onEncryptSubCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "algo",
						Aliases: []string{"a"},
						Value:   defaultAlgo,
						Usage:   "crypto algorithm which can be bcrypt, scrypt, argon2 or pbkdf2",
					},
				},
			},
		},
		Action: onGlobalCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func onGlobalCommand(ctx *cli.Context) (err error) {
	v := ctx.Bool("version")

	if v {
		return VersionAction(ctx)
	}

	return nil
}

func onCustomGenerateSubCommand(ctx *cli.Context) (err error) {
	// 1. process first argument
	len := ctx.Int("length")

	if ctx.NArg() > 0 {
		s := ctx.Args().Get(0)
		len, err = strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	// 2. process flags
	digits := ctx.Int("digits")
	symbols := ctx.Int("symbols")
	uppercase := ctx.Bool("uppercase")
	repeat := ctx.Bool("repeat")
	quiet := ctx.Bool("quiet")

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

func onGenerateSubCommand(ctx *cli.Context) (err error) {
	// 1. process first argument
	len := ctx.Int("length")

	if ctx.NArg() > 0 {
		s := ctx.Args().Get(0)
		len, err = strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	// 2. process flags
	digits := ctx.Bool("digits")
	symbols := ctx.Bool("symbols")
	uppercase := ctx.Bool("uppercase")
	repeat := ctx.Bool("repeat")
	quiet := ctx.Bool("quiet")

	if !quiet {
		fmt.Printf("Generating password %d characters long...\n", len)
	}

	digitsN := 0

	if digits {
		digitsN = defaultDigits
	}

	symbolsN := 0
	if symbols {
		symbolsN = defaultSymbols
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

func onEncryptSubCommand(ctx *cli.Context) (err error) {
	algo := ctx.String("algo")

	if ctx.NArg() > 0 {
		algo = ctx.Args().Get(0)
	}

	supportedHash := false

	switch algo {
	case algoBcrypt, algoScrypt, algoArgon2, algoPbkdf2:
		supportedHash = true
	}

	if !supportedHash {
		return fmt.Errorf("`%s` is not a supported hash or key derivation function", algo)
	}

	stdinPass, err := ask.HiddenAsk("New secure password: ")

	if err != nil {
		return err
	}

	stdinSalt, err := ask.HiddenAsk("A random salt (at least 8 bytes): ")

	if err != nil {
		return err
	}

	passwdEnc := hashPassword(stdinPass, stdinSalt, algo)

	fmt.Println("Password encrypted:")
	fmt.Println(passwdEnc)

	return nil
}

// borrowed from https://github.com/go-gitea/gitea/blob/master/models/user.go
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
