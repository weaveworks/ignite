package losetup

const (
	// general constants
	NameSize = 64
	KeySize  = 32
	Major    = 7

	// paths
	LoopControlPath = "/dev/loop-control"

	// loop flags
	FlagsReadOnly  = 1
	FlagsAutoClear = 4
	FlagsPartScan  = 8
	FlagsDirectIO  = 16

	// loop filter types
	CryptNone      = 0
	CryptXor       = 1
	CryptDes       = 2
	CryptFish2     = 3
	CryptBlow      = 4
	CryptCast128   = 5
	CryptIdea      = 6
	CryptDummy     = 9
	CryptSkipjack  = 10
	CryptCryptoApi = 18
	MaxCrypt       = 20

	// ioctl commands
	SetFd       = 0x4C00
	ClrFd       = 0x4C01
	SetStatus   = 0x4C02
	GetStatus   = 0x4C03
	SetStatus64 = 0x4C04
	GetStatus64 = 0x4C05
	ChangeFd    = 0x4C06
	SetCapacity = 0x4C07
	SetDirectIO = 0x4C08

	CtlAdd     = 0x4C80
	CtlRemove  = 0x4C81
	CtlGetFree = 0x4C82
)
