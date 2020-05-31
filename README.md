# Blox KeyVault


[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

Blox KeyVault is a library wrapping all major functionality an eth 2.0 validator will need:
  - [Multi key storage implementations](https://github.com/bloxapp/KeyVault/tree/master/stores)
  - [Signer](https://github.com/bloxapp/KeyVault/tree/master/validator_signer)
  - [Slashing protection](https://github.com/bloxapp/KeyVault/tree/master/slashing_protectors)
  - Tests


### Installation

 ```sh
go get github.com/bloxapp/KeyVault
   ```

### Tech

KeyVault uses a number of open source projects:

* [Prysm remote signer proto](github.com/wealdtech/eth2-signer-api) - Set of calls and structures for prysm's remote wallet
* [HashiCorp's Vault](https://github.com/hashicorp/vault) - As one of the key storage and ecnryption implementations


### Development and Debugging (via [goland](https://www.jetbrains.com/go/))

Want to contribute? Great!

The project uses Go modules for fast builds and clean dependencies.
Make a change in your file and instantaneously see your updates!

##### Install pre-requirement
* Install latest goland version
* clone the project:
```sh
git clone https://github.com/bloxapp/KeyVault.git
```

Basic use:
```go
	
// minimal configuration
store := stores.NewInMemStore()
encryptor := enc.NewPlainTextEncryptor()

options := vault.WalletOptions{}
options.SetEncryptor(encryptor)
options.SetStore(store)
options.SetWalletName("wallet")
options.SetWalletPassword("password")
options.EnableSimpleSigner(true) // false by default. if false, signer will not be available

// key management in one place
vault, _ := NewKeyVault(options)
account, _ := vault.wallet.CreateAccount("account","unlock_password")
account, _ = vault.wallet.AccountByName("account")

// manage all validator duty signatures
// pb package is [what used in prysm](github.com/wealdtech/eth2-signer-api/pb/v1)
vault.signer.SignBeaconProposal(*pb.SignBeaconProposalRequest)
vault.signer.SignBeaconAttestation(*pb.SignBeaconAttestationRequest)
```
