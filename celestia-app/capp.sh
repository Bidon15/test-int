#!/bin/bash

set -e

capp=/usr/bin/celestia-appd
# text="help"
# Display help and initialising later
CHAIN_ID="test"
KEY_TYPE="test"
COINS_TYPE="10000000000stake,100000000000celes"

$capp $text
$capp init $CHAIN_ID --chain-id $CHAIN_ID

# Creating the account  
$capp keys add $NODE_NAME --keyring-backend=$KEY_TYPE
node_addr=$($capp keys show $NODE_NAME -a --keyring-backend $KEY_TYPE)

$capp add-genesis-account $node_addr $COINS_TYPE --keyring-backend $KEY_TYPE
$capp gentx $NODE_NAME 5000000000stake --keyring-backend=$KEY_TYPE --chain-id $CHAIN_ID
$capp collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.celestia-app/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.celestia-app/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.celestia-app/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.celestia-app/config/config.toml
# Open up rest api
sed -i '104 s/enable = false/enable = true/' ~/.celestia-app/config/app.toml
sed -i '107 s/enable = false/enable = true/' ~/.celestia-app/config/app.toml

$capp start