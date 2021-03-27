#!/bin/sh
set -ex
sudo cp *.timer *.service /etc/systemd/system/.
sudo systemctl daemon-reload
for s in poold waterfall verifystate spafilter; do
    sudo systemctl disable $s || true # clear out any previous bogus configuration
    sudo systemctl enable $s
done
