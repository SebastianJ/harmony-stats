#!/usr/bin/env bash

echo "Installing Harmony Stats"
curl -LOs http://tools.harmony.one.s3.amazonaws.com/release/linux-x86_64/stats && chmod u+x stats
echo "Harmony Stats is now ready to use!"
echo "Invoke it using ./stats - see ./stats --help for all available options!"
