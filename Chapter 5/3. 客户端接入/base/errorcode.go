package base

const (
	ErrorCodeOK int32 = 1

	ErrorCodeServiceBusy = 100

	ErrorCodeInternalError = 101

	// --- Account Error --- //

	ErrorCodeAccountIdPasswordWrong int32 = 1000

	// ---  Login Error --- //

	ErrorCodeInvalidPlatformOrToken int32 = 2000

	// --- Gate Error --- //

	ErrorCodeInvalidGateToken int32 = 3000
	ErrorCodeClientBadRequest int32 = 3001

	// --- Game Error --- //

	ErrorInvalidSecret int32 = 4000
)
