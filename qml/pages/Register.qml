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
            acceptText: "Register"
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: qsTr("Connect with Signal")
        }

        TextField {
            id: telField
            width: parent.width
            inputMethodHints: Qt.ImhDialableCharactersOnly | Qt.ImhNoPredictiveText
            validator: RegExpValidator{ regExp: /[0-9]+/;}
            label: "Phone number (E.164 format)"
            placeholderText: "Your Phone Number"
            placeholderColor: Theme.highlightColor
            horizontalAlignment: TextInput.AlignLeft
            color: errorHighlight? "red" : Theme.primaryColor
            EnterKey.onClicked: parent.focus = true
        }

        TextArea {
            anchors.horizontalCenter: parent.horizontalCenter
            width: parent.width
            horizontalAlignment: TextEdit.Center
            readOnly: true
            text: qsTr("Signal will call you with a 6-digit verification code. Please be ready to write this down.")
        }

    }
}
