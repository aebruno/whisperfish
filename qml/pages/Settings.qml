
import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
	id: settingsPage
    property QtObject countryCodeCombo : countryCode
	SilicaFlickable {
		anchors.fill: parent
		contentWidth: parent.width
		contentHeight: col.height + Theme.paddingLarge

        PullDownMenu {
            MenuItem {
                text: qsTr("Refresh Contacts")
                onClicked: {
                    whisperfish.refreshContacts()
                    whisperfish.refreshSessions()
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
                text: whisperfish.phoneNumber()
            }
            TextArea {
                id: identity
                anchors.horizontalCenter: parent.horizontalCenter
                readOnly: true
                font.pixelSize: Theme.fontSizeSmall
                width: parent.width
                label: "Identity"
                text: whisperfish.identity()
            }
            SectionHeader {
                text: qsTr("General")
            }
            ValueButton {
                id: countryCode
                anchors.horizontalCenter: parent.horizontalCenter
                label: qsTr("Country Code")
                value: whisperfish.settings().countryCode
                onClicked: {
                    var cd = pageStack.push(Qt.resolvedUrl("CountryCodeDialog.qml"))
                    cd.setCountryCode.connect(function(code) {
                        value = code
                        whisperfish.settings().countryCode = code
                        whisperfish.saveSettings()
                    })
                }
            }
            TextSwitch {
                id: enableNotify
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Enable Notifications")
                checked: whisperfish.settings().enableNotify
                onCheckedChanged: {
                    whisperfish.settings().enableNotify = checked
                    whisperfish.saveSettings()
                }
            }
            TextSwitch {
                id: saveAttachments
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Save Attachments")
                checked: whisperfish.settings().saveAttachments
                onCheckedChanged: {
                    whisperfish.settings().saveAttachments = checked
                    whisperfish.saveSettings()
                }
            }
		}
	}
}
