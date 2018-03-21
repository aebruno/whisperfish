import QtQuick 2.2
import Sailfish.Silica 1.0

Dialog {
    id: phoneNumberDialog
    objectName: "phoneNumberDialog"
    property string tel
    signal selected(string tel)
    canAccept: isValid()

    function isValid() {
        if(telField.errorHighlight){
            return false
        }
        if(ContactModel.format(telField.text) == "") {
            return false
        }

        return true
    }

    onDone: {
        if (result == DialogResult.Accepted && !telField.errorHighlight) {
            tel = ContactModel.format(telField.text)
            phoneNumberDialog.selected(tel)
        }
    }

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            //: Enter phone number accept
            //% "Done"
            acceptText: qsTrId("whisperfish-new-message-accept-enter-number")
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
    }
}
