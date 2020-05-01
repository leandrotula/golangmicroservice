package environment

import "os"

const key = "AUTHORIZATION"

func RetrieveAuthorizationHeader() string  {

	return os.Getenv(key)

}
