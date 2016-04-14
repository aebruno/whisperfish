===============================================================================
Whisperfish - Signal client for Sailfish OS
===============================================================================

Whisperfish is a native `Signal <https://www.whispersystems.org/>`_ client for
`Sailfish OS <https://sailfishos.org/>`_. The code is based off the TextSecure
client for the Ubuntu Phone written by `janimo <https://github.com/janimo/textsecure-qml>`_. 

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

- [ ] Registration
- [ ] Contact Discovery
- [ ] Direct and group messages
- [ ] Photo/video attachments in both direct and group mode
- [ ] Storing conversations
- [ ] Encrypted message store

-------------------------------------------------------------------------------
Developing
-------------------------------------------------------------------------------

Whisperfish is written in go. First need to setup `MerSDK
<https://sailfishos.org/develop/sdk-overview/develop-installation-article/>`_
and install the go runtime. More details `here <https://github.com/nekrondev/jolla_go>`_.

Whisperfish uses a patched version of `go-qml <https://github.com/go-qml/qml>`_ 
for use with Safilish Silica UI. A complete patched version can be found 
`here <https://github.com/aebruno/qml/tree/sailfish>`_. If you followed the
jolla_go instructions above from nekrondev, you'll need to replace the 
~/src/gopkg.in/qml.v1 package with this version::

    $ cd ~/src/gopkg.in/
    $ rm -Rf qml.v1
    $ git clone https://github.com/aebruno/qml.git qml.v1
    $ cd qml.v1
    $ git checkout sailfish
    $ go install
    $ sb2 -O use-global-tmp -t SailfishOS-armv7hl ~/go/bin/linux_arm/go install

Additionally, the following go packages are required. To install them run::

    $ go get github.com/Sirupsen/logrus
    $ go get github.com/janimo/textsecure
    $ go get github.com/jmoiron/sqlx
    $ go get github.com/mattn/go-sqlite3
    $ go get gopkg.in/yaml.v2
    $ go get github.com/ttacon/libphonenumber

Clone the Whisperfish repo::

    $ git clone https://github.com/aebruno/whisperfish.git
    $ cd whisperfish
    $ go build
    $ mb2 build

If you have the SailfishOS Emulator you can install the rpm into the emulator
with directly with::

    $ ./deploy

-------------------------------------------------------------------------------
i18n Translations (help wanted)
-------------------------------------------------------------------------------

Whisperfish supports i18n translations. To translate the application strings in
your language run (for example German)::

    $ cd whisperfish
    $ sb2 lupdate qml/pages -ts qml/i18n/qml_de.ts
    [edit qml_de.ts]
    $ sb2 lrelease qml/i18n/qml_de.ts -qm qml/i18n/qml_de.qm

-------------------------------------------------------------------------------
License
-------------------------------------------------------------------------------

Copyright (C) 2016 Andrew E. Bruno

Whisperfish is free software: you can redistribute it and/or modify it under the
terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <http://www.gnu.org/licenses/>.
