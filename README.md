# gnark2zkv

Utilities for converting Groth16 artifacts from gnark format to zkVerify format.
The repository includes:

- A Go module `gnark2json` with a `SaveToJSON` function that can be used to serialize Gnark proofs, vks and public inputs in JSON format. An example of its usage can be found in `examples/main.go` .
- A command line utility `gnark2zkv` which takes the previously generated JSONs and convert them in a format suitable for submission to zkVerify

## Building
```
npm install
```

## Testing
```
npm test
```

## Usage
- Converting a gnark Groth16 zk-proof stored in file `proof.json`
    ```bash
    node gnark2zkv convert-proof proof.json -c bn128 # for BN128 curve
    ```
    ```bash
    node gnark2zkv convert-proof proof.json -c bls12381 # for BLS12-381 curve
    ```
- Converting a gnark Groth16 verification key stored in file `verification_key.json`
    ```bash
    node gnark2zkv convert-vk vk.json -c bn128 # for BN128 curve
    ```
     ```bash
    node gnark2zkv convert-vk vk.json -c bls12381 # for BLS12-381 curve
    ```
- Converting Groth16 public inputs stored in file `public.json`
    ```bash
    node gnark2zkv convert-public public.json -c bn128 # for BN128 curve
    ```
    ```bash
    node gnark2zkv convert-public public.json -c bls12381 # for BLS12-381 curve
    ```
By default these commands print their output to stdout, but it's possible redirect output to a file via the `-o` or `--out` option.
