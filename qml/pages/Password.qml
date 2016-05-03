import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: passwordDialog
    objectName: "passwordDialog"
    property string password

    canAccept: isValid()

    function isValid() {
        if(passwordField.errorHighlight){
            return false
        }
        if(!whisperfish.hasEncryptionKeys() && passwordField2.errorHighlight){
            return false
        }

        if(!whisperfish.hasEncryptionKeys() && passwordField.text != passwordField2.text){
            return false
        }

        return true
    }

    onDone: {
        if (result == DialogResult.Accepted && isValid()) {
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
            text: whisperfish.hasEncryptionKeys() ? qsTr("Enter your password") : qsTr("Set your password")
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

        TextField {
            id: passwordField2
            width: parent.width
            inputMethodHints: Qt.ImhNoPredictiveText
            visible: !whisperfish.hasEncryptionKeys()
            validator: RegExpValidator{ regExp: /.{6,}/;}
            label: "Verify Password"
            placeholderText: "Type password again"
            placeholderColor: Theme.highlightColor
            horizontalAlignment: TextInput.AlignLeft
            color: errorHighlight ? "red" : Theme.primaryColor
            EnterKey.onClicked: parent.focus = true
            echoMode: TextInput.Password
        }

        TextArea {
            anchors.horizontalCenter: parent.horizontalCenter
            visible: !whisperfish.hasEncryptionKeys()
            width: parent.width
            font.pixelSize: Theme.fontSizeTiny
            horizontalAlignment: TextEdit.Center
            readOnly: true
            text: qsTr("Whisperfish uses encrypted storage for identity keys and session state. This password is not stored anywhere and you will not be able to restore your data if you lose your password. IMPORTANT NOTE your messages are stored unencrypted on disk. Currently only identity keys and session state are stored encrypted.")
        }
    }
}
