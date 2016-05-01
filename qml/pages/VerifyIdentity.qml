import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: verifyIdentity
    objectName: "verifyIdentity"
    property string code

    canAccept: false

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            acceptText: ""
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: qsTr("Verify ") + messageModel.name
        }

        SectionHeader {
            text: qsTr("Their Identity (they read)")
        }

        TextArea {
            id: contactIdentity
            anchors.horizontalCenter: parent.horizontalCenter
            readOnly: true
            font.pixelSize: Theme.fontSizeSmall
            width: parent.width
            text: messageModel.identity
        }

        SectionHeader {
            text: qsTr("Your Identity (you read)")
        }

        TextArea {
            id: identity
            anchors.horizontalCenter: parent.horizontalCenter
            readOnly: true
            font.pixelSize: Theme.fontSizeSmall
            width: parent.width
            text: whisperfish.identity()
        }

    }
}
