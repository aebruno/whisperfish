===============================================================================
Whisperfish - Signal client for Sailfish OS
===============================================================================

Whisperfish is a native `Signal <https://www.whispersystems.org/>`_ client for
`Sailfish OS <https://sailfishos.org/>`_. Whisperfish uses the `Signal client
library for Go <https://github.com/aebruno/textsecure>`_ and `Qt binding for Go
<https://github.com/therecipe/qt>`_.  The user interface is heavily based on
the jolla-messages application written by Jolla Ltd.
  
-------------------------------------------------------------------------------
Project Status
-------------------------------------------------------------------------------

Whisperfish should be considered alpha software and used at your own risk. The
client is not an official Signal client and is not affiliated with Open Whisper
Systems. The code has not been audited by an expert in computer security or
cryptography. The goal of Whisperfish is to eventually become a stable, secure,
and robust Signal client for Sailfish OS. Code review and contributions are
welcome!

-------------------------------------------------------------------------------
Features
-------------------------------------------------------------------------------

- [x] Registration
- [x] Contact Discovery
- [x] Direct messages
- [x] Group messages
- [x] Storing conversations
- [x] Photo attachments
- [x] Encrypted identity and session store
- [x] Encrypted message store
- [x] Advanced user settings
- [ ] Multi-Device support (links with Signal Desktop)
- [ ] Encrypted local attachment store
- [ ] Archiving conversations

-------------------------------------------------------------------------------
Performance Tips
-------------------------------------------------------------------------------

Whisperfish connects to Signal using Websockets. For a better user experience
try adjusting the power settings on your Jolla to disable late suspend [1].
This should keep the network interfaces up and allow Whisperfish to maintain
websocket connections even when the device is in "sleep". This could
potentially impact your battery life depending on your usage. Otherwise
every time your device goes into deep sleep, the Websocket connection is broken
and you may not receive messages until the next time the OS wakes up and
Whisperfish reconnects.

To disable late suspend and enable "early suspend" run::

    $ mcetool --set-suspend-policy=early    

See here for more information.

1. https://together.jolla.com/question/55056/dynamic-pm-in-jolla/
2. http://talk.maemo.org/showpost.php?p=1401956&postcount=29
3. https://sailfishos.org/wiki/Sailfish_OS_Cheat_Sheet#Blocking_Device_Suspend

-------------------------------------------------------------------------------
Building from source
-------------------------------------------------------------------------------

*These instructions assume you're running Linux*

1. Install Go >= 1.10 and setup a proper `GOPATH <https://golang.org/doc/code.html#GOPATH>`_
   somewhere in your home directory, for example ``GOROOT=$HOME/projects/goroot/go`` and
   ``GOPATH=$HOME/projects/go``.

2. Install `Glide <https://glide.sh/>`_

3. Install Qt 5.7.0 the official `prebuilt package <https://download.qt.io/official_releases/qt/5.7/5.7.0/qt-opensource-linux-x64-android-5.7.0.run>`_

4. Install `Sailfish OS SDK <https://sailfishos.org/wiki/Application_SDK_Installation>`_ (version
   Beta-1801 or later). Ensure ``~/SailfishOS/vmshare/devices.xml`` exists. If not,
   run ``~/SailfishOS/bin/qtcreator`` once and it should create this file.

5. Clone whisperfish and download dependencies::

    $ git clone https://github.com/aebruno/whisperfish.git
    $ cd whisperfish
    $ glide install

6. Download and install Go QT bindings. This will clone
   https://github.com/therecipe/qt and checkout the specific version that has
   been known to work with whisperfish::

    $ ./build.sh bootstrap-qt

7. Run qtmoc and qtminimal (this is run on your local machine not the sdk)::

    $ ./build.sh prep

8. Create i386 and arm Go binaries for use in Sailfish SDK. Note this only
   needs to be done once::

    $ ./build.sh setup-sdk

9a.Login to SDK and setup environment. Installing Go in your home directory
   will provide SDK access to the Go binaries for compiling whisperfish.
   Note only needs to be done once::

    $ ssh -p 2222 -i ~/SailfishOS/vmshare/ssh/private_keys/engine/mersdk mersdk@localhost
    $ vim ~/.bashrc
    [add these lines to .bashrc]
    export GOROOT=/home/src1/projects/goroot/go
    export GOPATH=/home/src1/projects/go
    export PATH=$GOROOT/bin/linux_386:$PATH

    $ source ~/.bashrc


9b.Compile newer version of qemu and create new target for sb2. Older versions
   of qemu don't work well with Go::

    $ sudo zypper -n install libtool zlib-devel glib2-devel flex bison gcc pkgconfig glib2-static glibc-static make pcre-static
    $ cd $HOME
    $ mkdir src; cd src
    $ curl -O -L https://download.qemu.org/qemu-2.5.1.tar.bz2
    $ tar xjf qemu-2.5.1.tar.bz2
    $ ./configure --target-list=arm-softmmu,arm-linux-user
    $ make
    $ sudo make install
    $ sb2-init -L --sysroot=/ -C --sysroot=/ -c /usr/local/bin/qemu-arm -m sdk-build -n -N -t / SailfishOSgo-armv7hl /srv/mer/toolings/SailfishOS-2.1.4.13/opt/cross/bin/armv7hl-meego-linux-gnueabi-gcc


10. Login to SDK, compile whisperfish, and deploy to emulator::

    $ ssh -p 2222 -i ~/SailfishOS/vmshare/ssh/private_keys/engine/mersdk mersdk@localhost
    $ cd $GOPATH/src/github.com/aebruno/whisperfish
    $ ./build.sh compile
    $ ./build.sh i18n
    $ ./build.sh deploy

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
i18n Translations (help wanted)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Whisperfish supports i18n translations and uses Text ID Based Translations. See
`here <http://doc.qt.io/qt-5/linguist-id-based-i18n.html>`_ for more info. To
translate the application strings in your language run (for example German)::

    $ ssh -p 2222 -i ~/SailfishOS/vmshare/ssh/private_keys/engine/mersdk mersdk@localhost
    $ cd $GOPATH/src/github.com/aebruno/whisperfish
    $ sb2 lupdate qml/ -ts qml/i18n/whisperfish_de.ts
    [edit whisperfish_de.ts]
    $ sb2 lrelease -idbased qml/i18n/whisperfish_de.ts -qm qml/i18n/whisperfish_de.qm

Currently translations are only accepted through github pull requests.

-------------------------------------------------------------------------------
License
-------------------------------------------------------------------------------

Copyright (C) 2016-2018 Andrew E. Bruno

Whisperfish is free software: you can redistribute it and/or modify it under the
terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <http://www.gnu.org/licenses/>.
