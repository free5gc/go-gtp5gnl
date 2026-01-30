#!/bin/bash

make -C ../cmd/gogtp5g-link
make -C ../cmd/gogtp5g-tunnel

mv ../cmd/gogtp5g-link/gogtp5g-link .
mv ../cmd/gogtp5g-tunnel/gogtp5g-tunnel .