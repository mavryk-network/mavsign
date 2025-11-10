package preamble

import (
	// Install all backends
	_ "github.com/mavryk-network/mavsign/pkg/vault/aws"
	_ "github.com/mavryk-network/mavsign/pkg/vault/azure"
	_ "github.com/mavryk-network/mavsign/pkg/vault/cloudkms"
	_ "github.com/mavryk-network/mavsign/pkg/vault/file"
	_ "github.com/mavryk-network/mavsign/pkg/vault/hashicorp"
	_ "github.com/mavryk-network/mavsign/pkg/vault/ledger"
	_ "github.com/mavryk-network/mavsign/pkg/vault/mem"
	_ "github.com/mavryk-network/mavsign/pkg/vault/pkcs11"
	_ "github.com/mavryk-network/mavsign/pkg/vault/yubi"
)
