# go-gtp5gnl

go-gtp5gnl provides a netlink library about gtp5g for Go.

## License

This software is released under the Apache 2.0 License, see LICENSE

## Usage
### List all PDR/FAR/QER
```
# ./gtp5g-tunnel list [pdr/far/qer]
./gtp5g-tunnel list pdr
```
### Get/Del/Add/Mod PDR/FAR/QER
```
# ./gtp5g-tunnel [get/del/add/mod] [PDR/FAR/QER] [interface_name] [seid] [id] [option]
./gtp5g-tunnel add pdr upfgtp0 1 3 --pcd 99
```
- options
    ```
    PDR OPTIONS

            --pcd <precedence>

            --hdr-rm <outer-header-removal>

            --far-id <existed-far-id>

            --ue-ipv4 <pdi-ue-ipv4>

            --f-teid <i-teid> <local-gtpu-ipv4>

            --sdf-desp <description-string>

            --sdf-tos-traff-cls <tos-traffic-class>

            --sdf-scy-param-idx <security-param-idx>

            --sdf-flow-label <flow-label>

            --sdf-id <id>

            --qer-id <id>

    FAR OPTIONS

            --action <apply-action>

            --hdr-creation <description> <o-teid> <peer-ipv4> <peer-port>

    QER OPTIONS

            --qer-id <qer-id>

            --qfi-id <qfi-id> [Value range: {0..63}]

            --rqi-d <rqi> [Value range: {0=not triggered, 1=triggered}]

            --ppp <ppp> [Value range: {0=not present, 1=present}]

            --ppi <ppi> [Value range: {0..7}]
    ```