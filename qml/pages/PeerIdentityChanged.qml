import QtQuick 2.0
import Sailfish.Silica 1.0

Dialog {
    id: peerChangedDialog
    objectName: "peerChangedDialog"
    property var source

    onDone: {
        if (result == DialogResult.Accepted) {
            ClientWorker.resetPeerIdentity("yes")
        } else {
            ClientWorker.resetPeerIdentity("no")
        }
    }

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            //: Reset peer identity accept text
            //% "Confirm"
            acceptText: qsTrId("whisperfish-reset-peer-accept")
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            font.bold: true
            //: Peer identity not trusted 
            //% "Peer identity is not trusted"
            text: qsTrId("whisperfish-peer-not-trusted")
        }

        TextArea {
            anchors.horizontalCenter: parent.horizontalCenter
            width: parent.width
            horizontalAlignment: TextEdit.Center
            readOnly: true
            //: Peer identity not trusted message
            //% "WARNING: %1 identity is no longer trusted. Tap Confirm to reset peer identity."
            text: qsTrId("whisperfish-peer-not-trusted-message").arg(source)
        }

    }
}
