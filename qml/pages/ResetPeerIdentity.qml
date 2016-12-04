import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: resetPeerDialog
    objectName: "resetPeerDialog"
    property var source

    onDone: {
        if (result == DialogResult.Accepted) {
            Prompt.resetPeerIdentity("yes")
        } else {
            Prompt.resetPeerIdentity("no")
        }
    }

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            acceptText: "Confirm"
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            text: qsTr("Peer identity is not trusted")
        }

        TextArea {
            anchors.horizontalCenter: parent.horizontalCenter
            width: parent.width
            horizontalAlignment: TextEdit.Center
            readOnly: true
            text: qsTr("WARNING: "+source+" identity is no longer trusted. Tap Confirm to reset peer identity.")
        }

    }
}
