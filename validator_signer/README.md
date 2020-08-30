# Blox Eth Key Manager - Validator Signer


[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

Validator Signer has the responsibility to sign the basic 3 operations an eth 2.0 validator needs:

    - sign attestation
    - sign block proposal
    - sign attestation aggregation
    - return available public keys


### Instantiation

 ```golang
    signer := &signer.SimpleSigner{
		wallet,
		slashingProtector
	}
   ```
