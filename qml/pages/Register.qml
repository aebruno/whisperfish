import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: registerDialog
    objectName: "registerDialog"
    property string tel

    canAccept: !telField.errorHighlight

    onDone: {
        if (result == DialogResult.Accepted && !telField.errorHighlight) {
            tel = telField.text
            Prompt.phoneNumber(tel)
        }
    }

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            //: Register accept text
            //% "Register"
            acceptText: qsTrId("whisperfish-register-accept")
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            //: Registration message
            //% "Connect with Signal"
            text: qsTrId("whisperfish-registration-message")
        }

        TextField {
            id: telField
            width: parent.width
            inputMethodHints: Qt.ImhDialableCharactersOnly | Qt.ImhNoPredictiveText
            validator: RegExpValidator{ regExp: /[0-9]+/;}
            //: Phone number input
            //% "Phone number (E.164 format)"
            label: qsTrId("whisperfish-phone-number-input-label")
            //: Phone number placeholder
            //% "18875550100"
            placeholderText: qsTrId("whisperfish-phone-number-input-placeholder")
            placeholderColor: Theme.highlightColor
            horizontalAlignment: TextInput.AlignLeft
            color: errorHighlight? "red" : Theme.primaryColor
            EnterKey.onClicked: parent.focus = true
        }

        TextSwitch {
            id: shareContacts
            anchors.horizontalCenter: parent.horizontalCenter
            //: Share contacts label
            //% "Share Contacts"
            text: qsTrId("whisperfish-share-contacts-label")
            checked: SettingsBridge.boolValue("share_contacts")
            onCheckedChanged: {
                if(checked != SettingsBridge.boolValue("share_contacts")) {
                    SettingsBridge.boolSet("share_contacts", checked)
                }
            }
        }

        TextArea {
            anchors.horizontalCenter: parent.horizontalCenter
            width: parent.width
            horizontalAlignment: TextEdit.Center
            readOnly: true
            //: Registration directions
            //% "Signal will call you with a 6-digit verification code. Please be ready to write this down."
            text: qsTrId("whisperfish-registration-directions")
        }

    }
}
