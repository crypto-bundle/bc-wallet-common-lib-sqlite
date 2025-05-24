# bc-wallet-common-lib-sqlite

## Description 

Library for manage sqlite config and connection.

Library contains:
* common sqlite config struct
* connection manager
* small function-helpers for work with transaction statement

## Usage example

Examples of create connection and write database communication code

### Config and connection
```go
package main

import (
	commonEnvConfig "github.com/crypto-bundle/bc-wallet-common-lib-configs/pkg/envconfig"
  commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/sqlite"
)

func main() {
}


```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

* **bc-wallet-common-lib-sqlite** is licensed under the [MIT NON-AI](./LICENSE) License.