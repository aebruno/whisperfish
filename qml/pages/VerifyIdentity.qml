import QtQuick 2.2
import Sailfish.Silica 1.0

Page {
    id: verifyIdentity
    objectName: "verifyIdentity"

    SilicaFlickable {
        anchors.fill: parent
        contentHeight: column.height

        RemorsePopup { id: remorse }

        PullDownMenu {
            MenuItem {
                //: Reset secure session menu item
                //% "Reset Secure Session"
                text: qsTrId("whisperfish-reset-session-menu")
                enabled: MessageModel.peerIdentity.length > 0
                onClicked: {
                    //: Reset secure session remorse message
                    //% "Resetting secure session"
                    remorse.execute(qsTrId("whisperfish-reset-session-message"),
                        function() {
                            console.log("Resetting secure session: "+MessageModel.peerTel)
                            MessageModel.endSession(MessageModel.peerTel)
                        })
                }
            }
        }

        Column {
            id: column
            width: parent.width
            spacing: Theme.paddingLarge

            PageHeader {
                title: MessageModel.peerName
            }

            SectionHeader {
                //: Verify safety numbers
                //% "Verify safety numbers"
                text: qsTrId("whisperfish-verify-contact-identity-title")
            }

            TextArea {
                id: numericFingerprint
                horizontalAlignment: TextEdit.Center
                readOnly: true
                width: parent.width
                text: MessageModel.numericFingerprint(SetupWorker.localId, MessageModel.peerTel)
            }

            TextArea {
                id: fingerprintDirections
                anchors.horizontalCenter: parent.horizontalCenter
                readOnly: true
                font.pixelSize: Theme.fontSizeSmall
                width: parent.width
                //: Numeric fingerprint instructions
                //% "If you wish to verify the security of your end-to-end encryption with %1, compare the numbers above with the numbers on their device."
                text: qsTrId("whisperfish-numeric-fingerprint-directions").arg(MessageModel.peerName)
            }
        }
    }
}
