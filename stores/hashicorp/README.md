# Blox KeyVault - Hashicorp Vault


[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

Hashicorp Vault is an advanced, open source, secret management platform. It has a concept of "sealing" data and saving it in various locations.
Saving an the data encrypted and only "unsealing" it in runtime helps with securing the validator's keys.

Unsealing can be from a quorum of people, independently which makes it ideal for organizational use. 

### Tech

We use Hashicorp's sdk

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
account, _ = vault.wallet.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")

// manage all validator duty signatures
// pb package is [what used in prysm](github.com/wealdtech/eth2-signer-api/pb/v1)
vault.signer.SignBeaconProposal(*pb.SignBeaconProposalRequest)
vault.signer.SignBeaconAttestation(*pb.SignBeaconAttestationRequest)
```
