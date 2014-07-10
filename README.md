#Bitregister

### This project is still in active development and it is neither finished nor functional yet.

This project is an attempt to create a server to handle Bitcoin payments. 

The idea is to place this server as a middleware between the wallet (bitcoind) and the backend of the system. This middleware exposes a REST API via HTTP that the backend consumes to create orders and keep track of them.

The main motivation behind this project is to offer a self hosted Bitcoin payments processor as a measure to encourage decentralization, as opposed to most solutions offered today that relay on a third party for storing their clients' funds.