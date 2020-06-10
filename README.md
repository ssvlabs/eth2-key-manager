# Blox KeyVault
:hammer_and_wrench: UNDER HEAVY CONSTRUCTION :hammer_and_wrench:


[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

Blox KeyVault is a library wrapping all major functionality an eth 2.0 validator will need:
  - [Multi key storage implementations](https://github.com/bloxapp/KeyVault/tree/master/stores)
  - [Signer](https://github.com/bloxapp/KeyVault/tree/master/validator_signer)
  - [Slashing protection](https://github.com/bloxapp/KeyVault/tree/master/slashing_protectors)
  - [HD wallet](https://github.com/bloxapp/KeyVault/tree/master/wallet_hd) (EIP-2333,2334,2335 compliant)
  - Tests

### Installation

 ```sh
go get github.com/bloxapp/KeyVault
   ```

### Security and Architecture
KeyVault is built in an hierarchy starting with the concept of a [Portfolio](https://github.com/bloxapp/KeyVault/blob/master/core/portfolio.go) which represents the seed.<br/>
A portfolio can then create [wallets](https://github.com/bloxapp/KeyVault/blob/master/core/wallet.go) under itself and a wallet can create [accounts](https://github.com/bloxapp/KeyVault/blob/master/core/account.go) under itself.

An account is the entity that will ultimately signs transactions.<br/> 
Wallets and accounts are derived according to [EIP-2334](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md#validator-keys):<br/>
1) Withdrawal key: m/12381/3600/wallet_index/0<br/>
2) Validation key: m/12381/3600/wallet_index/0/account_index

The seed and private keys are never held in memory by one of the objects but rather we use our [DerivableKey](https://github.com/bloxapp/KeyVault/blob/master/core/derivable_key.go) object which asks the storage for the decrypted seed for each operation.<br/>
This is done so to not mistakenly print to console or serialize an object with the secret in plain text in it.

![](https://github.com/bloxapp/KeyVault/blob/master/images/arch_overview.png?raw=true)


Basic use:
```go
	
// minimal configuration
options := vault.WalletOptions{}
options.SetStore(stores.NewInMemStore())

// key management in one place
vault, _ := NewKeyVault(options)
wallet, _ := vault.CreateWallet("wallet")
validator, _ = wallet.CreateValidatorAccount("account")
withdrawal, _ = wallet.GetWithdrawalAccount() // only 1 per wallet
```
