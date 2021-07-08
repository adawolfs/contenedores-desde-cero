#!/usr/bin/env bash

skopeo --insecure-policy copy docker://alpine:latest oci:alpine:latest
skopeo --insecure-policy copy docker://centos:latest oci:centos:latest
skopeo --insecure-policy copy docker://supertest2014/nyan:latest oci:nyan:latest

umoci unpack --rootless --image alpine:latest alpine-bundle
umoci unpack --rootless --image centos:latest centos-bundle
umoci unpack --rootless --image nyan:latest nyan-bundle

mkdir -p containers
mv nyan-bundle/rootfs containers/nyan
mv centos-bundle/rootfs containers/centos
mv alpine-bundle/rootfs containers/alpine

rm -rf alpine*
rm -rf nyan*
rm -rf centos*
