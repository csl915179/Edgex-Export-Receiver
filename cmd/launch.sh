#!/bin/bash
#
# Copyright (c) 2018
# Tencent
#
# SPDX-License-Identifier: Apache-2.0
#

DIR=$PWD
CMD=cmd

# Kill all edgex-export-receiver* stuff
function cleanup {
	pkill eedgex-export-receiver
}

cd $CMD
exec -a edgex-export-receiver ./edgex-export-receiver &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done