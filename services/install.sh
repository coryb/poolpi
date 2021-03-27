#!/bin/sh
set -ex
sudo cp *.timer *.service /etc/systemd/system/.
sudo systemctl daemon-reload
for service in poold; do
    sudo systemctl disable $service || true # clear out any previous bogus configuration
    sudo systemctl enable $service
done

for timer in waterfall verifystate spafilter; do
    # clear out any previous bogus configuration
    sudo systemctl disable $timer || true
    sudo systemctl disable $timer.timer || true

    sudo systemctl enable $timer.timer
done
