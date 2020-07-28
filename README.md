# Blox KeyVault
:hammer_and_wrench: UNDER HEAVY CONSTRUCTION :hammer_and_wrench:


[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

Blox KeyVault is a library wrapping all major functionality an eth 2.0 validator will need:
  - [Multi storage implementations](https://github.com/bloxapp/KeyVault/tree/master/stores)
  - [Signer](https://github.com/bloxapp/KeyVault/tree/master/validator_signer)
  - [Slashing protection](https://github.com/bloxapp/KeyVault/tree/master/slashing_protection)
  - [HD wallet](https://github.com/bloxapp/KeyVault/tree/master/wallet_hd) (EIP-2333,2334,2335 compliant)
  - Tests

### Installation

 ```sh
go get github.com/bloxapp/KeyVault
   ```

### Security and Architecture
KeyVault is the entry point to manage all operations, in it sits a unique wallet and accounts.<br/> 
KeyVault <- Wallet <- [Accounts]


An account is the entity that ultimately signs transactions.<br/> 
Wallets and accounts are derived according to [EIP-2334](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md#validator-keys):<br/>
1) Withdrawal key: m/12381/3600/account_index/0<br/>
2) Validation key: m/12381/3600/account_index/0/0<br/>

The seed is needed just to execute specific operations like creating new accounts or signing with the withdrawal key. <br/><br/>

Examples:
- [Basic Use]()