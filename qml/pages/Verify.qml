import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: verifyDialog
    objectName: "verifyDialog"
    property string code

    canAccept: !codeField.errorHighlight

    onDone: {
        if (result == DialogResult.Accepted && !codeField.errorHighlight) {
            code = codeField.text
            codeEntered(code)
        }
    }

    signal codeEntered(string text)

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            acceptText: "Verify"
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: qsTr("Verify Device")
        }

        TextField {
            id: codeField
            width: parent.width
            inputMethodHints: Qt.ImhDigitsOnly | Qt.ImhNoPredictiveText
            validator: RegExpValidator{ regExp: /[0-9]+/;}
            label: "Code"
            placeholderText: "Code"
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
            text: qsTr("Signal will now send you a confirmation code via SMS message. Please enter it here.")
        }

    }
}
