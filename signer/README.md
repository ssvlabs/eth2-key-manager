# Eth Key Manager - Validator Signer


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
