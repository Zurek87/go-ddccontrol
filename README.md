# go-ddccontrol
Go lib to use ddccontrol.

Gtk3 Icon to change brightness.

![Go GTK!](https://github.com/zurek87/go-ddccontrol/raw/master/static/pic1.png "Menu icon!")



# Install ddccontrol

[ddccontrol](https://github.com/ddccontrol/ddccontrol)

Best way is install from github.

After installation check existence of ```ddccontrol.pc``` in pkgconfig directory.


# Hints:

To run without sudo:
```
# groupadd i2c
# usermod -aG i2c <myusername>
# echo 'KERNEL=="i2c-[0-9]*", GROUP="i2c"' >> /etc/udev/rules.d/10-local_i2c_group.rules
```