# Example

In this document, it will cover using `gogtp5g-link` to setup up a gtp5g network device, using `gogtp5g-tunnel` to add PDU session's PDR/FAR in the device, and also show how to use `gogtp5g-tunnel` get the list of rules.

## Buld gogtp5g-link & gogtp5g-tunnel

```bash
./build.sh
```

This script will build gogtp5g-link and go gtp5g-tunnel bin file from cmd directory and move them to this example directory for later using.

## Case 1: Add One PDU Session

1. Setup script

    For adding one basic pdu session, we need to add a set of downlink and uplink PDR/FAR.

    The `ran.sh` will bring up a gtp5g network device and add the PDRs and FARs.

    ```bash
    ./ran.sh
    ```

2. Check PDR/FAR

    To check if the PDR and FAR is added successfully, run `gogtp5g-tunnel` with `list`.

    ```bash
    ./gogtp5g-tunnel list pdr
    ```

    ```bash
    ./gogtp5g-tunnel list far
    ```

3. Cleanup

    After checking, the cleanup scrip can help you to cleanup the environment.

    ```bash
    ./cleanup.sh
    ```

## Case 2: Add Two PDU Session

For setting two PDU session (to set of PDRs/FARs), we can just re-run the operation in `ran.sh`.

We've provide a script `ran-dual-pdu.sh` to insert two PDRs/FARs in the gtp5g device.

```bash
./ran-dual-pdu.sh
```

## Conclusion

With the above demonstration, user can write the customized script for generate the corrsponding environment.
