package config

import (
	"testing"

	"github.com/mavryk-network/mavbingo/v2/b58"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavsign/pkg/hashmap"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type testCase struct {
	title               string
	src                 string
	expect              *Config
	expectParseError    string
	expectValidateError string
}

func mustPKH(src string) crypt.PublicKeyHash {
	pkh, err := b58.ParsePublicKeyHash([]byte(src))
	if err != nil {
		panic(err)
	}
	return pkh
}

var testCases = []testCase{
	{
		title: "Valid",
		src: `---
base_dir: $HOME/.mavsign
server:
  address: :6732
  utility_address: :9583

vaults:
  kms:
    driver: cloudkms
    config:
      project: mavsign
      location: europe-north1
      key_ring: hsm-ring

mavryk:
  mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt:

  mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm:
    log_payloads: true
    allow:
      generic:
        - transaction
        - endorsement
      block:
      endorsement:
`,
		expect: &Config{
			BaseDir: "$HOME/.mavsign",
			Server: ServerConfig{
				Address:        ":6732",
				UtilityAddress: ":9583",
			},
			Vaults: map[string]*VaultConfig{
				"kms": {
					Driver: "cloudkms",
					Config: yaml.Node{
						Kind: 4,
						Tag:  "!!map",
						Content: []*yaml.Node{
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "project",
								Line:   11,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "mavsign",
								Line:   11,
								Column: 16,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "location",
								Line:   12,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "europe-north1",
								Line:   12,
								Column: 17,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "key_ring",
								Line:   13,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "hsm-ring",
								Line:   13,
								Column: 17,
							},
						},
						Line:   11,
						Column: 7,
					},
				},
			},
			Mavryk: hashmap.NewPublicKeyHashMap([]hashmap.PublicKeyKV[*MavrykPolicy]{
				{
					Key: mustPKH("mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt"),
					Val: nil,
				},
				{
					Key: mustPKH("mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm"),
					Val: &MavrykPolicy{
						LogPayloads: true,
						Allow: map[string][]string{
							"generic":     {"transaction", "endorsement"},
							"block":       nil,
							"endorsement": nil,
						},
					},
				},
			}),
		},
	},
	{
		title: "InvalidBase58",
		src: `---
base_dir: $HOME/.mavsign
server:
  address: :6732
  utility_address: :9583

vaults:
  kms:
    driver: cloudkms
    config:
      project: mavsign
      location: europe-north1
      key_ring: hsm-ring

mavryk:
  111111111111111111111111111111111111:
`,
		expectParseError: "mavbingo: base58Check decoding error: invalid checksum",
	},
	{
		title: "InvalidType",
		src: `---
base_dir: $HOME/.mavsign
server:
  address: :6732
  utility_address: :9583

vaults:
  kms:
    driver: cloudkms
    config:
      project: mavsign
      location: europe-north1
      key_ring: hsm-ring

mavryk:
  edpkv45regue1bWtuHnCgLU8xWKLwa9qRqv4gimgJKro4LSc3C5VjV:
`,
		expectParseError: "mavbingo: unknown public key prefix",
	},
	{
		title: "NoBaseDir",
		src: `---
server:
  address: :6732
  utility_address: :9583

vaults:
  kms:
    driver: cloudkms
    config:
      project: mavsign
      location: europe-north1
      key_ring: hsm-ring

mavryk:
  mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt:

  mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm:
    log_payloads: true
    allow:
      generic:
        - transaction
        - endorsement
      block:
      endorsement:
`,
		expect: &Config{
			Server: ServerConfig{
				Address:        ":6732",
				UtilityAddress: ":9583",
			},
			Vaults: map[string]*VaultConfig{
				"kms": {
					Driver: "cloudkms",
					Config: yaml.Node{
						Kind: 4,
						Tag:  "!!map",
						Content: []*yaml.Node{
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "project",
								Line:   10,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "mavsign",
								Line:   10,
								Column: 16,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "location",
								Line:   11,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "europe-north1",
								Line:   11,
								Column: 17,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "key_ring",
								Line:   12,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "hsm-ring",
								Line:   12,
								Column: 17,
							},
						},
						Line:   10,
						Column: 7,
					},
				},
			},
			Mavryk: hashmap.NewPublicKeyHashMap([]hashmap.PublicKeyKV[*MavrykPolicy]{
				{
					Key: mustPKH("mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt"),
					Val: nil,
				},
				{
					Key: mustPKH("mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm"),
					Val: &MavrykPolicy{
						LogPayloads: true,
						Allow: map[string][]string{
							"generic":     {"transaction", "endorsement"},
							"block":       nil,
							"endorsement": nil,
						},
					},
				},
			}),
		},
		expectValidateError: "Key: 'Config.BaseDir' Error:Field validation for 'BaseDir' failed on the 'required' tag",
	},
	{
		title: "InvalidAddress",
		src: `---
base_dir: $HOME/.mavsign
server:
  address: xxxx
  utility_address: :9583

vaults:
  kms:
    driver: cloudkms
    config:
      project: mavsign
      location: europe-north1
      key_ring: hsm-ring

mavryk:
  mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt:

  mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm:
    log_payloads: true
    allow:
      generic:
        - transaction
        - endorsement
      block:
      endorsement:
`,
		expect: &Config{
			BaseDir: "$HOME/.mavsign",
			Server: ServerConfig{
				Address:        "xxxx",
				UtilityAddress: ":9583",
			},
			Vaults: map[string]*VaultConfig{
				"kms": {
					Driver: "cloudkms",
					Config: yaml.Node{
						Kind: 4,
						Tag:  "!!map",
						Content: []*yaml.Node{
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "project",
								Line:   11,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "mavsign",
								Line:   11,
								Column: 16,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "location",
								Line:   12,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "europe-north1",
								Line:   12,
								Column: 17,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "key_ring",
								Line:   13,
								Column: 7,
							},
							{
								Kind:   8,
								Tag:    "!!str",
								Value:  "hsm-ring",
								Line:   13,
								Column: 17,
							},
						},
						Line:   11,
						Column: 7,
					},
				},
			},
			Mavryk: hashmap.NewPublicKeyHashMap([]hashmap.PublicKeyKV[*MavrykPolicy]{
				{
					Key: mustPKH("mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt"),
					Val: nil,
				},
				{
					Key: mustPKH("mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm"),
					Val: &MavrykPolicy{
						LogPayloads: true,
						Allow: map[string][]string{
							"generic":     {"transaction", "endorsement"},
							"block":       nil,
							"endorsement": nil,
						},
					},
				},
			}),
		},
		expectValidateError: "Key: 'Config.Server.Address' Error:Field validation for 'Address' failed on the 'hostname_port' tag",
	},
	{
		title: "EmptyVaultData",
		src: `---
base_dir: $HOME/.mavsign
server:
  address: :6732
  utility_address: :9583

vaults:
  kms:

mavryk:
  mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt:

  mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm:
    log_payloads: true
    allow:
      generic:
        - transaction
        - endorsement
      block:
      endorsement:
`,
		expect: &Config{
			BaseDir: "$HOME/.mavsign",
			Server: ServerConfig{
				Address:        ":6732",
				UtilityAddress: ":9583",
			},
			Vaults: map[string]*VaultConfig{
				"kms": nil,
			},
			Mavryk: hashmap.NewPublicKeyHashMap([]hashmap.PublicKeyKV[*MavrykPolicy]{
				{
					Key: mustPKH("mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt"),
					Val: nil,
				},
				{
					Key: mustPKH("mv3CXEtMtTwnnroNp3fxkWsBV1M5xP9uMCmm"),
					Val: &MavrykPolicy{
						LogPayloads: true,
						Allow: map[string][]string{
							"generic":     {"transaction", "endorsement"},
							"block":       nil,
							"endorsement": nil,
						},
					},
				},
			}),
		},
		expectValidateError: "Key: 'Config.Vaults[kms]' Error:Field validation for 'Vaults[kms]' failed on the 'required' tag",
	},
}

func TestConfig(t *testing.T) {
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			var result Config
			err := yaml.Unmarshal([]byte(test.src), &result)
			if test.expectParseError != "" {
				require.EqualError(t, err, test.expectParseError)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expect, &result)
				err := Validator().Struct(&result)
				if test.expectValidateError != "" {
					require.EqualError(t, err, test.expectValidateError)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}
