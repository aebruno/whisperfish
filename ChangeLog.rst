===============================================================================
ChangeLog
===============================================================================

v0.4.0 (unreleased)
---------------------------

* Built using SailfishOSSDK-Beta-1611-Qt5-linux-64
* Tested on Sailfish OS 2.0.5.6 (Haapajoki)
* Major code refactor to use new Go QT bindings
* Viewing conversations now use native QAbstractList models which should
  improve performance
* Fix #45 The attachment directory is now configurable and can be changed to a
  location searched by the gallery
* Fix #6 and #57 Notifications no longer replace. There is a new notification
  for each message
* Fix #58 Incognito mode should be working again
* Fix #55 (partially) Add command line option for manually
  encrypting/decrypting database
* Add option to disable sharing contacts with Signal
* Fix #52 Enable quick scroll
* Add ability to resend messages
* Fix #63 Add support for numeric fingerprints
* Add CLI tool for adding extensions to attachment file names

`v0.3.0`_ (2016-06-07)
---------------------------

* Built using SailfishOSSDK-Beta-1602-Qt5-linux-64
* Tested on Sailfish OS 2.0.1.11 (Taalojärvi)
* Fix #40 Add sound/LED to notifications
* Fix #35 copy to clipboard

`v0.2.0`_ (2016-06-06)
---------------------------

* Second alpha release
* Built using SailfishOSSDK-Beta-1602-Qt5-linux-64
* Tested on Sailfish OS 2.0.1.11 (Taalojärvi)
* Fix #32 Keyboard closes when message arrives in active conversation bug 
* Fix #9 Screen flickering
* Fix #25 Send button doesn't stay in place
* Fix #28 Remove cover action main page
* Fix #36 Fix incognito mode cancel

`v0.1.1`_ (2016-05-14)
---------------------------

* First alpha release
* Built using SailfishOSSDK-Beta-1511-Qt5-linux-64
* Tested on Sailfish OS 2.0.0.10 (Saimaa)

.. _v0.1.1: https://github.com/aebruno/whisperfish/releases/tag/v0.1.1
.. _v0.2.0: https://github.com/aebruno/whisperfish/releases/tag/v0.2.0
.. _v0.3.0: https://github.com/aebruno/whisperfish/releases/tag/v0.3.0
