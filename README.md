# bitcoin-prober

This is a simple program that sends a `version` message to a given `host:port` and outputs details about the remote BCH/BTC node.

Here are some example outputs:

```
$ ./bitcoin-prober --help
Usage of ./bitcoin-prober:
  -address string
    	Address to probe
  -network string
    	Network (BCH or BTC) (default "BCH")
  -verbose
    	Be verbose

$ ./bitcoin-prober --address seed.bchd.cash
Resolved seed.bchd.cash to 70.174.210.133:8333
Probing 70.174.210.133:8333 on the BCH network...
70.174.210.133 is located in US, Phoenix
UserAgent: /Bitcoin ABC:0.19.7(EB32.0)/
Services: SFNodeNetwork|SFNodeBloom|SFNodeBitcoinCash|SFNodeNetworkLimited
ProtocolVersion: 70015
LastBlock: 590811
RelayTx: true

$ ./bitcoin-prober --network BTC --address 141.105.69.133:8333
Probing 141.105.69.133:8333 on the BTC network...
141.105.69.133 is located in RU
UserAgent: /Satoshi:0.18.0/
Services: SFNodeNetwork|SFNodeBloom|SFNodeWitness|SFNodeNetworkLimited
ProtocolVersion: 70015
LastBlock: 584885
RelayTx: true
```

_This was quickly hacked together, improvements are welcome!_

## License

AGPLv3+
