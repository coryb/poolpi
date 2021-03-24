## About
This is project to simulate a remote control for the Hayward Prologic pool "automation" system.

There is a similar better supported Python project here: https://github.com/swilson/aqualogic
Also inspired by some excellent wire sleuthing: http://www.desert-home.com/p/swimming-pool.html

The basic components:
* Raspberry Pi (using 3B+, any model should work)
* RS485 <=> Serial UART converter
  * Used https://smile.amazon.com/gp/product/B010723BCE but there are many similar parts that will probably all work the same

## Wiring Diagram
```
RPi Pin #4 (5V)  -> UART VCC
RPi Pin #6 (GND) -> UART GND
RPi Pin #8 (TX)  -> UART TX
RPi Pin #10 (RX) -> UART RX

RS485 (A+)  -> Prologic Remote COMM Pin 2 (BLK)
RS485 (B-)  -> Prologic Remote COMM Pin 3 (YEL)
RS485 (GND) -> Prologic Remote COMM Pin 4 (GRN)
```

## Rasberry Pi Setup
We just need the serial port enabled.  This can be done by setting `enable_uart=1` in `/boot/config.txt` and rebooting.  If setup correctly you should see something like:
```
$ ls -l /dev/serial0
lrwxrwxrwx 1 root root 5 Mar 21 18:17 /dev/serial0 -> ttyS0
```

## Disclaimer
This project is not affiliated with or endorsed by Hayward Industries Inc. in any way.