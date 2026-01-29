#!/bin/bash
source ./config

sudo ip addr add ${UE_IP} dev lo >> /dev/null 2>&1
sudo ip addr add ${UE_IP2} dev lo >> /dev/null 2>&1

sudo killall -9 gogtp5g-link
sudo ./gogtp5g-link add gtp5gtest --ran &
sleep 0.2

sudo ./gtp5g-tunnel add far gtp5gtest 1 --action 2
sudo ./gtp5g-tunnel add far gtp5gtest 2 --action 2 --hdr-creation 0 78 ${UPF_IP} 2152

sudo ./gtp5g-tunnel add pdr gtp5gtest 1 --pcd 1 --hdr-rm 0 --ue-ipv4 ${UE_IP} --f-teid 87 ${RAN_IP} --far-id 1
sudo ./gtp5g-tunnel add pdr gtp5gtest 2 --pcd 2 --ue-ipv4 ${UE_IP} --far-id 2

sudo ./gtp5g-tunnel add far gtp5gtest 3 --action 2
sudo ./gtp5g-tunnel add far gtp5gtest 4 --action 2 --hdr-creation 0 79 ${UPF_IP} 2152

sudo ./gtp5g-tunnel add pdr gtp5gtest 3 --pcd 1 --hdr-rm 0 --ue-ipv4 ${UE_IP2} --f-teid 88 ${RAN_IP} --far-id 3
sudo ./gtp5g-tunnel add pdr gtp5gtest 4 --pcd 2 --ue-ipv4 ${UE_IP2} --far-id 4

sudo ip r add ${DN_CIDR} dev gtp5gtest
sudo ip r add ${DN_CIDR2} dev gtp5gtest

