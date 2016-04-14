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
            numberEntered(tel)
            console.log(tel)
        }
    }

    signal numberEntered(string text)

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            acceptText: "Register"
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: qsTr("Your phone number")
        }


        TextField {
            id: telField
            width: parent.width
            inputMethodHints: Qt.ImhDialableCharactersOnly | Qt.ImhNoPredictiveText
            validator: RegExpValidator{ regExp: /[0-9]{11,11}/;}
            label: "Phone number (E.164 format)"
            placeholderText: "14155552671"
            placeholderColor: Theme.highlightColor
            horizontalAlignment: TextInput.AlignLeft
            color: errorHighlight? "red" : Theme.primaryColor
            EnterKey.onClicked: parent.focus = true
        }

    }
}
