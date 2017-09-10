#!/bin/bash -e
cp source/Dockerfile output/.
cp releases/cf-mgmt-linux output/cf-mgmt
chmod +x output/cf-mgmt
