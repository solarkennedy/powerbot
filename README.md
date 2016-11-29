## powerbot

[![Build Status](https://travis-ci.org/solarkennedy/powerbot.svg?branch=master)](https://travis-ci.org/solarkennedy/powerbot)

`powerbot` is an IRC bot designed to interface with the
[digi-rc-switch](https://github.com/solarkennedy/digi-rc-switch) usb
peripheral. This peripheral can send RF commands to remote controllable
outlets.

## Config

See the example [config file](https://github.com/solarkennedy/powerbot/blob/master/powerbot.yaml)
for how to configure `powerbot`.

Once configured, it can send arbitrary codes or respond to pre-programmed aliases:

```
   @kyle | powerbot: christmas off
powerbot | Sent out code 95500 for christmas
   @kyle | powerbot: christmas on
powerbot | Sent out code 95491 for christmas
   @kyle | powerbot: code 1234
powerbot | Sent out code 1234
```

## Compatible devices

The following devices are known to be controllable with this code:

- [Etekcity Wireless Remote Control Electrical Outlet Switch for Household Appliances](http://www.etekcity.com/product/100068.html), receiver model BH9938U

## TODO

- [ ] Multi-channel support
