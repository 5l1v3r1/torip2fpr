# torip2fpr
Given Tor relay IP addresses, find their corresponding fingerprints.  The
fingerprints are extracted from
[CollecTor's](https://collector.torproject.org/archive/relay-descriptors/consensuses/)
archives, which you need in order to run the tool.

# Usage
    $ torip2fpr -datadir /path/to/collector/consensuses/ -addrfile /path/to/addresses

The IP address file has to have the following format:

    label_1: addr_x, addr_x+1, ...
    label_2: addr_y, addr_y+1, ...
    ...

The output is CSV-formatted.  Every line contains `fingerprint`, `address`,
`label`.

# Contact
Contact: Philipp Winter <phw@nymity.ch>  
OpenPGP fingerprint: `B369 E7A2 18FE CEAD EB96  8C73 CF70 89E3 D7FD C0D0`
