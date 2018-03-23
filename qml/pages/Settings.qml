
import QtQuick 2.2
import Sailfish.Silica 1.0

Page {
    id: settingsPage
    property QtObject countryCodeCombo : countryCode

    RemorsePopup {
        id: remorse
        onCanceled: {
            incognitoModeSwitch.checked = !incognitoModeSwitch.checked
        }
    }

    SilicaFlickable {
        anchors.fill: parent
        contentWidth: parent.width
        contentHeight: col.height + Theme.paddingLarge

        PullDownMenu {
            MenuItem {
                //: Linked devices menu option
                //% "Linked Devices"
                text: qsTrId("whisperfish-settings-linked-devices-menu")
                onClicked: pageStack.push(Qt.resolvedUrl("LinkedDevices.qml"))
            }
            MenuItem {
                //: Reconnect menu
                //% "Reconnect"
                text: qsTrId("whisperfish-settings-reconnect-menu")
                onClicked: {
                    ClientWorker.reconnect()
                }
            }
            MenuItem {
                //: Refresh contacts menu
                //% "Refresh Contacts"
                text: qsTrId("whisperfish-settings-refresh-contacts-menu")
                onClicked: {
                    ContactModel.refresh()
                    contactList.refresh()
                    SessionModel.reload()
                }
            }
        }

        VerticalScrollDecorator {}

        Column {
            id: col
            spacing: Theme.paddingLarge
            width: parent.width
            PageHeader {
                //: Settings page title
                //% "Whisperfish Settings"
                title: qsTrId("whisperfish-settings-title")
            }
            SectionHeader {
                //: Settings page My identity section label
                //% "My Identity"
                text: qsTrId("whisperfish-settings-identity-section-label")
            }
            TextField {
                id: phone
                anchors.horizontalCenter: parent.horizontalCenter
                readOnly: true
                width: parent.width
                //: Settings page My phone number
                //% "Phone"
                label: qsTrId("whisperfish-settings-my-phone-number")
                text: SetupWorker.phoneNumber
            }
            TextArea {
                id: identity
                anchors.horizontalCenter: parent.horizontalCenter
                readOnly: true
                font.pixelSize: Theme.fontSizeSmall
                width: parent.width
                //: Settings page Identity label
                //% "Identity"
                label: qsTrId("whisperfish-settings-identity-label")
                text: SetupWorker.identity
            }
            SectionHeader {
                //: Settings page notifications section
                //% "Notifications"
                text: qsTrId("whisperfish-settings-notifications-section")
            }
            TextSwitch {
                id: enableNotify
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page notifications enable
                //% "Enabled"
                text: qsTrId("whisperfish-settings-notifications-enable")
                checked: SettingsBridge.boolValue("enable_notify")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("enable_notify")) {
                        SettingsBridge.boolSet("enable_notify", checked)
                    }
                }
            }
            TextSwitch {
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page notifications show message body
                //% "Show Message Body"
                text: qsTrId("whisperfish-settings-notifications-show-body")
                checked: SettingsBridge.boolValue("show_notify_message")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("show_notify_message")) {
                        SettingsBridge.boolSet("show_notify_message", checked)
                    }
                }
            }
            SectionHeader {
                //: Settings page general section
                //% "General"
                text: qsTrId("whisperfish-settings-general-section")
            }
            ValueButton {
                id: countryCode
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page country code
                //% "Country Code"
                label: qsTrId("whisperfish-settings-country-code")
                value: SettingsBridge.stringValue("country_code")
                onClicked: {
                    var cd = pageStack.push(Qt.resolvedUrl("CountryCodeDialog.qml"))
                    cd.setCountryCode.connect(function(code) {
                        value = code
                        SettingsBridge.stringSet("country_code", code)
                    })
                }
            }
            TextSwitch {
                id: saveAttachments
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page save attachments
                //% "Save Attachments"
                text: qsTrId("whisperfish-settings-save-attachments")
                checked: SettingsBridge.boolValue("save_attachments")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("save_attachments")) {
                        SettingsBridge.boolSet("save_attachments", checked)
                    }
                }
            }
            TextSwitch {
                id: shareContacts
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page share contacts
                //% "Share Contacts"
                text: qsTrId("Share Contacts")
                checked: SettingsBridge.boolValue("share_contacts")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("share_contacts")) {
                        SettingsBridge.boolSet("share_contacts", checked)
                    }
                }
            }
            TextSwitch {
                id: enableEnterSend
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page enable enter send
                //% "EnterKey Send"
                text: qsTrId("whisperfish-settings-enable-enter-send")
                checked: SettingsBridge.boolValue("enable_enter_send")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("enable_enter_send")) {
                        SettingsBridge.boolSet("enable_enter_send", checked)
                    }
                }
            }
            SectionHeader {
                //: Settings page advanced section
                //% "Advanced"
                text: qsTrId("whisperfish-settings-advanced-section")
            }
            TextSwitch {
                id: incognitoModeSwitch
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page incognito mode
                //% "Incognito Mode"
                text: qsTrId("whisperfish-settings-incognito-mode")
                checked: SettingsBridge.boolValue("incognito")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("incognito")) {
                        remorse.execute(
                            //: Restart whisperfish message
                            //% "Restart Whisperfish..."
                            qsTrId("whisperfish-settings-restarting-message"),
                            function() {
                                SettingsBridge.boolSet("incognito", checked)
                                SetupWorker.restart()
                        })
                    }
                }
            }
            TextSwitch {
                id: scaleImageAttachments
                anchors.horizontalCenter: parent.horizontalCenter
                //: Settings page scale image attachments
                //% "Scale JPEG Attachments"
                text: qsTrId("whisperfish-settings-scale-image-attachments")
                checked: SettingsBridge.boolValue("scale_image_attachments")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("scale_image_attachments")) {
                        SettingsBridge.boolSet("scale_image_attachments", checked)
                    }
                }
            }
            SectionHeader {
                //: Settings page stats section
                //% "Statistics"
                text: qsTrId("whisperfish-settings-stats-section")
            }
            DetailItem {
                //: Settings page websocket status
                //% "Websocket Status"
                label: qsTrId("whisperfish-settings-websocket")
                value: ClientWorker.connected ? 
                    //: Settings page connected message
                    //% "Connected"
                    qsTrId("whisperfish-settings-connected") : 
                    //: Settings page disconnected message
                    //% "Disconnected"
                    qsTrId("whisperfish-settings-disconnected")
            }
            DetailItem {
                //: Settings page unsent messages
                //% "Unsent Messages"
                label: qsTrId("whisperfish-settings-unsent-messages")
                value: MessageModel.unsentCount()
            }
            DetailItem {
                //: Settings page total sessions
                //% "Total Sessions"
                label: qsTrId("whisperfish-settings-total-sessions")
                value: SessionModel.count()
            }
            DetailItem {
                //: Settings page total messages
                //% "Total Messages"
                label: qsTrId("whisperfish-settings-total-messages")
                value: MessageModel.total()
            }
            DetailItem {
                //: Settings page total signal contacts
                //% "Signal Contacts"
                label: qsTrId("whisperfish-settings-total-contacts")
                value: ContactModel.total()
            }
            DetailItem {
                //: Settings page encrypted key store
                //% "Encrypted Key Store"
                label: qsTrId("whisperfish-settings-encrypted-keystore")
                value: SetupWorker.encryptedKeystore ? 
                    //: Settings page encrypted key store enabled
                    //% "Enabled"
                    qsTrId("whisperfish-settings-encrypted-keystore-enabled") : 
                    //: Settings page encrypted key store disabled
                    //% "Disabled"
                    qsTrId("whisperfish-settings-encrypted-keystore-disabled")
            }
            DetailItem {
                //: Settings page encrypted database
                //% "Encrypted Database"
                label: qsTrId("whisperfish-settings-encrypted-db")
                value: SettingsBridge.boolValue("encrypt_database") ? 
                    //: Settings page encrypted db enabled
                    //% "Enabled"
                    qsTrId("whisperfish-settings-encrypted-db-enabled") : 
                    //: Settings page encrypted db disabled
                    //% "Disabled"
                    qsTrId("whisperfish-settings-encrypted-db-disabled")
            }
        }
    }
}
