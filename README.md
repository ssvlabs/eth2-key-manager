# Blox Eth Key Manager

<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbloxapp%2Feth2-key-manager.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbloxapp%2Feth2-key-manager?ref=badge_shield) -->

Blox Eth Key Manager is a library wrapping all major functionality an eth 2.0 validator will need:
  - [Multi storage implementations](https://github.com/ssvlabs/eth2-key-manager/tree/master/stores)
  - [Signer](https://github.com/ssvlabs/eth2-key-manager/tree/master/validator_signer)
  - [Slashing protection](https://github.com/ssvlabs/eth2-key-manager/tree/master/slashing_protection)
  - [HD wallet](https://github.com/ssvlabs/eth2-key-manager/tree/master/wallet_hd) (EIP-2333,2334,2335 compliant)
  - Tests

### Installation

 ```sh
go get github.com/ssvlabs/eth2-key-manager
   ```

### Security and Architecture
`eth2keymanager` is the entry point to manage all operations, in it sits a unique wallet and accounts.<br/> 
`eth2keymanager` <- Wallet <- [Accounts]


An account is the entity that ultimately signs transactions.<br/> 
Wallets and accounts are derived according to [EIP-2334](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md#validator-keys):<br/>
1) Withdrawal key: m/12381/3600/account_index/0<br/>
2) Validation key: m/12381/3600/account_index/0/0<br/>

The seed is needed just to execute specific operations like creating new accounts or signing with the withdrawal key. <br/><br/>

Examples:
- [Basic Use]()
