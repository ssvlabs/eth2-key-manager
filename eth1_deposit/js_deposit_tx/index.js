var Web3 = require('web3');
const EthereumTx = require('ethereumjs-tx').Transaction;

var web3 = new Web3('https://goerli.infura.io/v3/d03b92aa81864faeb158166231b7f895');

const DepositContractAddress = "0x0F0F0fc0530007361933EaB5DB97d09aCDD6C1c8"; // TODO - should not be hard coded
const DepositContractABI = require("./eth2deposit.json");

const From = "<your addre>";
const FromPriv = Buffer.from(
    '<address priv key>',
    'hex',);


var depositContract = new web3.eth.Contract(DepositContractABI,DepositContractAddress);

// Get the params from KeyVault DepositData function
methodAbi = depositContract.methods.deposit(
    "", // pubkey
    "", // WithdrawalCredentials
    "", // signature
    ""  // depositDataRoot
).encodeABI();



web3
    .eth
    .getTransactionCount(From)
    .then(nonce => {
        const txParams = {
            nonce: "0x" + parseInt(nonce).toString(16),
            gasPrice: "0x59682F02", // 1500000002 wei // TODO - should not be hard coded
            gasLimit: '0x7A120', // 500K // TODO - should not be hard coded
            to: DepositContractAddress,
            value: web3.utils.numberToHex(web3.utils.toWei("32","ether")), // TODO - should not be hard coded
            data: methodAbi,
        }
        const tx = new EthereumTx(txParams, { chain: 'goerli', hardfork: 'petersburg' }); // TODO - should not be hard coded
        tx.sign(FromPriv);
        const raw = '0x' + tx.serialize().toString('hex');

        return web3.eth.sendSignedTransaction(raw);
    })
    .then (res => {
        console.log(res)
    })
    .catch (err => {
        console.log(err)
    })
