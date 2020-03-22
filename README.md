# rPasswd

> A secure random password generator tool.

This is small CLI tool to generate [secure random passwords](https://golang.org/pkg/crypto/rand/) as described by [AgileBits 1Password](https://discussions.agilebits.com/discussion/23842/how-random-are-the-generated-passwords).

## Install

```sh
go get -u github.com/joseluisq/rpasswd
```

Release binaries also available on [joseluisq/rpasswd/releases](https://github.com/joseluisq/rpasswd/releases)

## Usage

### Default mode

Default mode enables lower-upper cases, digits and symbols characters by default.

```
$ rpasswd gen --quiet
~}@3)&`?[]>/.#7^!:<=92",_\8{(+-16*%$|054⏎

$ rpasswd gen --quiet 64
#K!{Q.@\p?)3N+6:7<2xv8uV[O_ar"`w$}f1>G4=&0HT^],P~*/I-%(|Ce9i5JyE⏎

$ rpasswd gen --symbols=false --quiet
uPAZS1J0Qa69IDxUNq2i8E3CFo4Wv5jth7OfMseB⏎

$ rpasswd gen --symbols=false --digits=false --quiet
GREkMnyzTcVUIFCaOqsXembSPwrfAjHLJuBYDiZp⏎

$ rpasswd gen --symbols=false --uppercase=false -q 36
gfyihq8caomube3vz7n9k402dwts1px6l5jr⏎
```

### Custom mode

```
$ rpasswd --digits=0 --symbols=0 --quiet 40
ODxJAeMtUhkWFgQBToRqIdviNaHlyCSzGKPpXwuY⏎

$ rpasswd --digits=10 --symbols=0 --quiet 40
0jg587sOf9VmndquiC3Gt2XvQYRKHZw6UN1BD4JT⏎

$ rpasswd --digits=10 --symbols=10 --quiet 40
2Oj.7h:0?FG8U%W&4k|QT1-mt[i3a<lVB69ps}5C⏎

$ rpasswd --digits=10 --symbols=0 --repeat --quiet 40
1XZ8lvkJ74al6zqqCek8ars2V0hgB9NzeriJb3YE⏎

$ rpasswd --digits=10 --symbols=0 --repeat --quiet 64
sbF2ZJre1bNfVOnwhWhXf1iJ5ly9IBTSyHMhLQH427p6Y5MEpodAJmyXKymipYlk⏎
```

## Global options

```
$ rpasswd -h

NAME:
   rPasswd - A secure random password generator tool

USAGE:
   rpasswd [global options] command [command options] rpasswd [global options] <length>

   Example: rpasswd --symbols 20 --uppercase 32

DESCRIPTION:
   CLI tool to generate passwords as described by AgileBits 1Password[1]. The algorithm is commonly used when generating website passwords.

    [1] https://discussions.agilebits.com/discussion/23842/how-random-are-the-generated-passwords

COMMANDS:
   gen      Generate a random password including lower-upper cases, digits and symbols characters by default.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --digits value, -d value   quantity of digits (max. 10) (default: 0)
   --length value, -l value   Password length (default: 40)
   --quiet, -q                disallow verbose mode (default: false)
   --repeat, -r               allow characters to repeat (default: false)
   --symbols value, -s value  quantity of symbols (max. 30) (default: 0)
   --uppercase, -u            allow upper and lower cases (default: true)
   --help, -h                 show help (default: false)
```

## Command options

```
$ rpasswd gen -h

NAME:
   rpasswd gen - Generate a random password including lower-upper cases, digits and symbols characters by default.

USAGE:
   rpasswd gen [command options] <length>

OPTIONS:
   --length value, -l value  password length (default: 40)
   --digits, -d              enable 10 digit characters (default: true)
   --symbols, -s             enable 30 symbol characters (default: true)
   --uppercase, -u           allow upper-lower case characters. false for lowercase only (default: true)
   --repeat, -r              allow characters to repeat (default: false)
   --quiet, -q               disallow verbose mode (default: false)
   --help, -h                show help (default: false)
```

## Encryption

```sh
enc -h
NAME:
   main enc - Encrypt a given password using a key derivation function hash.

USAGE:
   main enc [command options] <hash function>

OPTIONS:
   --hash value, -s value  key derivation function hash like bcrypt, scrypt, argon2 or pbkdf2 (default: "pbkdf2")
   --help, -h              show help (default: false)
```

## Contributions

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in current work by you, as defined in the Apache-2.0 license, shall be dual licensed as described below, without any additional terms or conditions.

Feel free to send some [Pull request](https://github.com/joseluisq/rpasswd/pulls) or [issue](https://github.com/joseluisq/rpasswd/issues).

## License

This work is primarily distributed under the terms of both the [MIT license](LICENSE-MIT) and the [Apache License (Version 2.0)](LICENSE-APACHE).

© 2020 [Jose Quintana](https://git.io/joseluisq)
