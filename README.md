# Pando-eth-rpc-adaptor

The `Pando-eth-rpc-adaptor` project is aiming to provide an adaptor which translates the Pando RPC interface to the Ethereum RPC interface. 

## Setup

First, install **Go ** and set environment variables `GOPATH` , `GOBIN`, and `PATH`. Next, clone the Pando blockchain repo and install Pando following the steps below:

```
mkdir -p usr/local/go/src/github.com/pandoprojects 
cd usr/local/go/src/github.com/pandoprojects
git clone https://
cd usr/local/go/src/github.com/pandoprojects/pando
export GO111MODULE=on
make install
```

Next, clone the `pando-eth-rpc-adaptor` repo:

```
cd usr/local/go/src/github.com/pandoprojects
git clone https://github.com/pandoprojects/Pando-eth-rpc-libary.git pando-eth-rpc-adaptor
```

## Build and Install

### Build the binary under macOS or Linux
Following the steps below to build the `pando-eth-rpc-adaptor` binary and copy it into your `$GOPATH/bin`.

```
echo 'export PANDO_RPC=/usr/local/go/src/github.com/pandoprojects/pando-eth-rpc-adaptor' >> ~/.profile
echo 'export PANDO_RPC=/usr/local/go/src/github.com/pandoprojects/pando-eth-rpc-adaptor' >> ~/.bashrc
source ~/.bashrc && source ~/.profile
cd $PANDO_RPC
export GO111MODULE=on
make install
```

### Cross compilation for Windows
On a macOS machine, the following command should build the `pando-eth-rpc-adaptor.exe` binary under `build/windows/`

```
make windows
```

## Run the Adaptor with a local Pando private testnet

First, run a private testnet Pando node with its RPC port opened at 16888:

```
cd $PANDO_HOME
cp -r ./integration/pandoproject ../pandoproject
mkdir ~/.pandocli
cp -r ./integration/pandoproject/pandocli/* ~/.pandocli/
chmod 700 ~/.pandocli/keys/encrypted

pando start --config=../pandoprojects/node_eth_rpc 
choose a password 
```

Then, open another terminal, create the config folder for the RPC adaptor

```
mkdir -p ../pandoprojects/eth-rpc-adaptor
```

Use your favorite editor to open file `../pandoprojects/eth-rpc-adaptor/config.yaml`, paste in the follow content, save and close the file:

```
pando
  rpcEndpoint: "http://127.0.0.1:16888/rpc"
rpc:
  enabled: true
  httpAddress: "127.0.0.1"
  httpPort: 18888
  wsAddress: "127.0.0.1"
  wsPort: 18889
  timeoutSecs: 600 
  maxConnections: 2048
log:
  levels: "*:debug"
```

Then, launch the adaptor binary with the following command:

```
cd $PANDO_RPC
pando-eth-rpc-adaptor start --config=../pandoproject/eth-rpc-adaptor
```

The RPC adaptor will first create 10 test wallets, which will be useful for running tests with dev tools like Truffle, Hardhat. After the test wallets are created, the ETH RPC APIs will be ready for use.

## RPC APIs

The RPC APIs should conform to the Ethereum JSON RPC API standard: https://eth.wiki/json-rpc/API. We currently support the following Ethereum RPC APIs:

```
eth_chainId
eth_syncing
eth_accounts
eth_protocolVersion
eth_getBlockByHash
eth_getBlockByNumber
eth_blockNumber
eth_getUncleByBlockHashAndIndex
eth_getTransactionByHash
eth_getTransactionByBlockNumberAndIndex
eth_getTransactionByBlockHashAndIndex
eth_getBlockTransactionCountByHash
eth_getTransactionReceipt
eth_getBalance
eth_getStorageAt
eth_getCode
eth_getTransactionCount
eth_getLogs
eth_getBlockTransactionCountByNumber
eth_call
eth_gasPrice
eth_estimateGas
eth_sendRawTransaction
eth_sendTransaction
net_version
web3_clientVersion
```

The following examples demonstrate how to interact with the RPC APIs using the `curl` command:

```
# Query Chain ID
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":67}' http://localhost:18888/rpc

# Query synchronization status
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' http://localhost:18888/rpc

# Query block number
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://localhost:18888/rpc

# Query account PTX balance (should return an integer which represents the current PTX balance in wei)
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x2E833968E5bB786Ae419c4d13189fB081Cc43bab", "latest"],"id":1}' http://localhost:18888/rpc
```

Further detail of our smartcontract documentaiton please visit our [official Documentation site](https://docs.pandoproject.org/pandoproject/smart-contracts)
