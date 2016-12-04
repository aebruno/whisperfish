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
        if(!Backend.registered && passwordField2.errorHighlight){
            return false
        }

        if(!Backend.registered && passwordField.text != passwordField2.text){
            return false
        }

        return true
    }

    onDone: {
        if (result == DialogResult.Accepted && isValid()) {
            password = passwordField.text
            Prompt.password(password)
        }
    }

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader { }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: Backend.registered ? qsTr("Enter your password") : qsTr("Set your password")
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
            visible: !Backend.registered
            validator: RegExpValidator{ regExp: /.{6,}/;}
            label: "Verify Password"
            placeholderText: "Verify Password"
            placeholderColor: Theme.highlightColor
            horizontalAlignment: TextInput.AlignLeft
            color: errorHighlight ? "red" : Theme.primaryColor
            EnterKey.onClicked: parent.focus = true
            echoMode: TextInput.Password
        }

        TextArea {
            anchors.horizontalCenter: parent.horizontalCenter
            visible: !Backend.registered
            width: parent.width
            font.pixelSize: Theme.fontSizeTiny
            horizontalAlignment: TextEdit.Center
            readOnly: true
            text: qsTr("Whisperfish stores identity keys, session state, and local message data encrypted on disk. The password you set is not stored anywhere and you will not be able to restore your data if you lose your password. Note: Attachments are currently stored unencrypted. You can disable storing of attachments in the Settings page.")
        }
    }
}
