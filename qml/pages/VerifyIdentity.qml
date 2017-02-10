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
            //: Verify contact identity
            //% "Verify %1"
            text: qsTrId("whisperfish-verify-contact-identity-title").arg(MessageModel.peerName)
        }

        SectionHeader {
            //: Contact identity message
            //% "Their Identity (they read)"
            text: qsTrId("whisperfish-contact-identity-section")
        }

        TextArea {
            id: contactIdentity
            anchors.horizontalCenter: parent.horizontalCenter
            readOnly: true
            font.pixelSize: Theme.fontSizeSmall
            width: parent.width
            text: MessageModel.peerIdentity
        }

        SectionHeader {
            //: Your identity message
            //% "Your Identity (you read)"
            text: qsTrId("whisperfish-your-identity-section")
        }

        TextArea {
            id: identity
            anchors.horizontalCenter: parent.horizontalCenter
            readOnly: true
            font.pixelSize: Theme.fontSizeSmall
            width: parent.width
            text: SetupWorker.identity
        }

    }
}
