
import QtQuick 2.0
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
                text: qsTr("Linked Devices")
                onClicked: pageStack.push(Qt.resolvedUrl("LinkedDevices.qml"))
            }
            MenuItem {
                text: qsTr("Reconnect")
                onClicked: {
                    Backend.reconnect()
                }
            }
            MenuItem {
                text: qsTr("Refresh Contacts")
                onClicked: {
                    Backend.contactRefresh()
                }
            }
        }

        VerticalScrollDecorator {}

        Column {
            id: col
            spacing: Theme.paddingLarge
            width: parent.width
            PageHeader {
                title: qsTr("Whisperfish Settings")
            }
            SectionHeader {
                text: qsTr("My Identity")
            }
            TextField {
                id: phone
                anchors.horizontalCenter: parent.horizontalCenter
                readOnly: true
                width: parent.width
                label: "Phone"
                text: Backend.phoneNumber()
            }
            TextArea {
                id: identity
                anchors.horizontalCenter: parent.horizontalCenter
                readOnly: true
                font.pixelSize: Theme.fontSizeSmall
                width: parent.width
                label: "Identity"
                text: Backend.identity()
            }
            SectionHeader {
                text: qsTr("Notifications")
            }
            TextSwitch {
                id: enableNotify
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Enable")
                checked: SettingsBridge.boolValue("enable_notify")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("enable_notify")) {
                        SettingsBridge.boolSet("enable_notify", checked)
                    }
                }
            }
            TextSwitch {
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Show Message Body")
                checked: SettingsBridge.boolValue("show_notify_message")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("show_notify_message")) {
                        SettingsBridge.boolSet("show_notify_message", checked)
                    }
                }
            }
            SectionHeader {
                text: qsTr("General")
            }
            ValueButton {
                id: countryCode
                anchors.horizontalCenter: parent.horizontalCenter
                label: qsTr("Country Code")
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
                text: qsTr("Save Attachments")
                checked: SettingsBridge.boolValue("save_attachments")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("save_attachments")) {
                        SettingsBridge.boolSet("save_attachments", checked)
                    }
                }
            }
            SectionHeader {
                text: qsTr("Advanced")
            }
            TextSwitch {
                id: incognitoModeSwitch
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Incognito Mode")
                checked: SettingsBridge.boolValue("incognito")
                onCheckedChanged: {
                    if(checked != SettingsBridge.boolValue("incognito")) {
                        remorse.execute(
                            qsTr("Restarting whisperfish..."),
                            function() {
                                SettingsBridge.boolSet("incognito", checked)
                                Backend.restart()
                        })
                    }
                }
            }
            SectionHeader {
                text: qsTr("Statistics")
            }
            DetailItem {
                label: qsTr("Network Status")
                value: Backend.connected ? "Connected" : "Disconnected"
            }
            DetailItem {
                label: qsTr("Unsent Messages")
                value: MessageModel.unsentCount()
            }
            DetailItem {
                label: qsTr("Total Sessions")
                value: SessionModel.count()
            }
            DetailItem {
                label: qsTr("Total Messages")
                value: MessageModel.total()
            }
            DetailItem {
                label: qsTr("Signal Contacts")
                value: Backend.contactCount()
            }
            DetailItem {
                label: qsTr("Encrypted Key Store")
                value: Backend.hasEncryptedKeystore() ? qsTr("Enabled") : qsTr("Disabled")
            }
            DetailItem {
                label: qsTr("Encrypted Database")
                value: SettingsBridge.boolValue("encrypt_database") ? qsTr("Enabled") : qsTr("Disabled")
            }
        }
    }
}
