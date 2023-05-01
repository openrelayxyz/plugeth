# PluGeth

PluGeth is a fork of the [Go Ethereum Client](https://github.com/ethereum/go-ethereum)
(Geth) that implements a plugin architecture, allowing developers to extend
Geth's  capabilities in a number of different ways using plugins, rather than
having to create additional, new forks of Geth.

Documentation can be found [here](https://plugeth.org).

## Design Goals

The upstream Geth client exists primarily to serve as a client for the Ethereum
mainnet, though it also supports a number of popular testnets. Supporting the
Ethereum mainnet is a big enough challenge in its own right that the Geth team
generally avoids changes to support other networks, or to provide features only
a small handful of users would be interested in.

The result is that many projects have forked Geth. Some implement their own
consensus protocols or alter the behavior of the EVM to support other networks.
Others are designed to extract information from the Ethereum mainnet in ways
the standard Geth client does not support.

PluGeth aims to provide a single Geth fork that developers can choose to extend
rather than forking the Geth project. Out of the box, PluGeth behaves exactly
like upstream Geth, but by installing plugins written in Golang, developers can
extend its functionality in a wide variety of way.

### Contact Us

If you're trying to do something that isn't supported by the current plugin system, Reach out to us on [Discord](https://discord.gg/Epf7b7Gr) and we'll help you figure out how to make it work.

## System Requirements

System requirements will vary depending on which network you are connecting to.
On the Ethereum mainnet, you should have at least 8 GB RAM, 2 CPUs, and 350 GB
of SSD disks.

PluGeth relies on Golang's Plugin implementation, which is only supported on
Linux, FreeBSD, and macOS. Windows support is unlikely to be added in the
foreseeable future.

# Licensing Considerations

The Geth codebase is licensed under the LGPL. By linking with Geth, you have an
obligation to enable anyone you provide your plugin binaries to run against
their own modified versions of Geth. Because of how Golang plugins work
running against updated versions of Geth may require recompiling the plugin.

If you plan to license your plugin under the LGPL or a more permissive license,
you should be able to meet these requirements. If you plan to use your plugin
privately without distributing it, you should be fine. If you plan to release
your plugin without making the source available, you may find yourself in
violation of Geth's license unless you can provide a way to relink it against
more recent versions of Geth.


