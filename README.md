<p align="center">
  <a href="https://sourceforge.net/projects/desktopify/">
    <img src="https://a.fsdn.com/allura/p/desktopify/icon?1649542343&w=90">
  </a>
</p>

## Desktopify
**Ubuntu AppImage installer**

**Description**

A utility that attempts to reduce complexity of AppImage software execution. Instead of manually running the AppImage, Desktopify creates shortcuts and places the executable in an installation folder.

**Installation**
```sh
wget desktopify.tgz https://github.com/profkittyface/desktopify/releases/download/1.0/desktopify-1.0.tgz
tar xvf desktopify-1.0.tgz
sudo mv desktopify /usr/local/bin/desktopify
```

**Usage**
```sh
desktopify list
desktopify install Software.AppImage
desktopify remove Software.AppImage
```
