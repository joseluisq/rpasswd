package cmd

// supported hashing functions
const (
	algoBcrypt = "bcrypt"
	algoScrypt = "scrypt"
)

// supported key derivation hash functions
const (
	algoArgon2 = "argon2"
	algoPbkdf2 = "pbkdf2"
)

// default application values
const (
	// defaultDigits contains default quantity of digits available
	defaultDigits = 10
	// defaultSymbols contains default quantity of symbols available
	defaultSymbols = 30
	// defaultLength contains default characters length
	defaultLength = 40
	// defaultAlgo contains default algorithm
	defaultAlgo = algoPbkdf2
)

// application version values
var (
	versionNumber string = "devel"
	buildTime     string
)
