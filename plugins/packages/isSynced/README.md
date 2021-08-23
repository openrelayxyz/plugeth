## IsSynced

Returns an object with data about sync status of a node. The plugin was written as an extension of the native JSON-RPC method eth_syncing. Unlike eth_syncing however the call will always return an object. Along with data about the sync status the object contains boolean values corresponding to the presence of active peers and the ultimate status of the sync. 

Example usage. 

```json 

curl 127.0.0.1:8545 -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"plugeth_isSynced","params":[],"id":22}'
