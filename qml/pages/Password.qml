import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: passwordDialog
    objectName: "passwordDialog"
    property string password

    canAccept: !passwordField.errorHighlight

    onDone: {
        if (result == DialogResult.Accepted && !passwordField.errorHighlight) {
            password = passwordField.text
            passwordEntered(password)
        }
    }

    signal passwordEntered(string text)

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader { }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: qsTr("Enter storage password")
        }

        TextField {
            id: passwordField
            width: parent.width
            inputMethodHints: Qt.ImhNoPredictiveText
            validator: RegExpValidator{ regExp: /.{6,}/;}
            label: "Password"
            placeholderText: "Password"
            placeholderColor: Theme.highlightColor
            horizontalAlignment: TextInput.AlignLeft
            color: errorHighlight? "red" : Theme.primaryColor
            EnterKey.onClicked: parent.focus = true
            echoMode: TextInput.Password
        }

    }
}
