import QtQuick 2.2
import Sailfish.Silica 1.0
import "../delegates"

Page {
    id: conversation
    objectName: "conversation"
    property bool editorFocus
    onStatusChanged: {
        if(status == PageStatus.Active) {
            pageStack.pushAttached(Qt.resolvedUrl("VerifyIdentity.qml"))
        }
    }

    MessagesView {
        id: messages
        focus: true
        anchors.fill: parent

        model: MessageModel

        // Use a placeholder for the ChatTextInput to avoid re-creating the input
        header: Item {
            width: messages.width
            height: textInput.height
        }

        Column {
            id: headerArea
            y: messages.headerItem.y
            parent: messages.contentItem
            width: parent.width

            ChatTextInput {
                id: textInput
                width: parent.width
                contactName: MessageModel.peerName
                enabled: true
                editorFocus: conversation.editorFocus

                onSendMessage: {
                    var sid = MessageModel.createMessage(MessageModel.peerTel, text, "", attachmentPath, true)
                    if(sid > 0) {
                        // Update session model
                        SessionModel.add(sid, true)
                    }
                }
            }
        }
    }
}
