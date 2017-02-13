===============================================================================
Whisperfish - Signal client for Sailfish OS
===============================================================================

Whisperfish is a native `Signal <https://www.whispersystems.org/>`_ client for
`Sailfish OS <https://sailfishos.org/>`_. Whisperfish uses the `Signal client
library for Go <https://github.com/janimo/textsecure>`_ and `Qt binding for Go
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
- [ ] Encrypted attachments
- [ ] Archiving conversations

-------------------------------------------------------------------------------
Building from source
-------------------------------------------------------------------------------

*These instructions assume you're running Linux*

1. Install Go >= 1.7.1 and setup a proper `GOPATH <https://golang.org/doc/code.html#GOPATH>`_ 
   somewhere in your home directory, for example ``$HOME/projects/go``.

2. Install `Glide <https://glide.sh/>`_

3. Install Qt 5.7.0 the offical `prebuilt package <https://download.qt.io/official_releases/qt/5.7/5.7.0/qt-opensource-linux-x64-android-5.7.0.run>`_

4. Install `Sailfish OS SDK <https://sailfishos.org/wiki/Application_SDK_Installation>`_

5. Clone whisperfish and download dependencies::

    $ git clone https://github.com/aebruno/whisperfish.git
    $ cd whisperfish
    $ glide install

6. Download and install Go QT bindings. This will clone
   https://github.com/therecipe/qt and checkout the specific version that has
   been known to work with whisperfish::

    $ ./build.sh bootstrap-qt

7. Run qtmoc and qtminimal (these are run your local machine not the mersdk)::

    $ ./build.sh prep

8. Create i386 Go binary for use in Sailfish SDK. Note this only needs to be
   done once::

    $ ./build.sh setup-mer

9. Login to mersdk and setup environment. Note only needs to be done once::

    $ ssh -p 2222 -i ~/SailfishOS/vmshare/ssh/private_keys/engine/mersdk mersdk@localhost
    $ vim ~/.bashrc
    [add these lines to .bashrc]
    export GOROOT=/home/src1/projects/goroot/go
    export GOPATH=/home/src1/projects/go
    export PATH=$GOROOT/bin/linux_386:$PATH

    $ source ~/.bashrc

10. Login to mersdk, compile whisperfish, and deploy to emulator::

    $ ssh -p 2222 -i ~/SailfishOS/vmshare/ssh/private_keys/engine/mersdk mersdk@localhost
    $ cd $GOPATH/src/github.com/aebruno/whisperfish
    $ ./build.sh compile
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

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Making new releases
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

*Note*: The latest tag from the current git branch is used in the package
version (``mb2 -x``). To add git hashes to the package version modify the
``/usr/bin/mb2`` script with the following patch::

    --- mb2.orig    2016-05-19 02:44:04.015412275 +0000
    +++ /usr/bin/mb2        2016-05-19 02:44:13.722084593 +0000
    @@ -154,7 +154,7 @@
     fix_package_version() {
         [[ ! $OPT_FIX_VERSION ]] && return
     
    -    local tag=$(git describe --tags --abbrev=0 2>/dev/null)
    +    local tag=$(git describe --long --tags --dirty --always 2>/dev/null)
         if [[ -n $tag ]]; then
             # tagver piece copied from tar_git service
             if [[ $(echo $tag | grep "/") ]] ; then

-------------------------------------------------------------------------------
License
-------------------------------------------------------------------------------

Copyright (C) 2016-2017 Andrew E. Bruno

Whisperfish is free software: you can redistribute it and/or modify it under the
terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <http://www.gnu.org/licenses/>.
