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

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Building from source
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

*TODO* These docs are incomplete.

Whisperfish uses `Glide <https://glide.sh/>`_ for package management. Ensure you
have a working Go runtime and that your GOPATH is set. Whisperfish uses QT
bindings for Go. More information on installing this library can be found 
`here <https://github.com/therecipe/qt>`_. The build scripts assume
you've installed QT here: $HOME/Qt5.x.x. You'll also need the `MerSDK
<https://sailfishos.org/wiki/Application_SDK_Installation>`_.

To build Whisperfish from source::

    $ git clone https://github.com/aebruno/whisperfish.git
    $ cd whisperfish
    $ glide install
    $ ./build.sh bootstrap-qt
    $ ./build.sh prep

    $ ssh to merdsk
    $ cd $GOPATH/src/github.com/aebruno/whisperfish
    $ ./build.sh compile
    $ ./build.sh deploy


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
i18n Translations (help wanted)
-------------------------------------------------------------------------------

Whisperfish supports i18n translations and uses Text ID Based Translations. See
`here <http://doc.qt.io/qt-5/linguist-id-based-i18n.html>`_ for more info. To
translate the application strings in your language run (for example German)::

    $ cd whisperfish
    $ sb2 lupdate qml/ -ts qml/i18n/whisperfish_de.ts
    [edit whisperfish_de.ts]
    $ sb2 lrelease -idbased qml/i18n/whisperfish_de.ts -qm qml/i18n/whisperfish_de.qm

Currently translations are only accepted through github pull requests.

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
