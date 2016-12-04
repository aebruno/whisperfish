import QtQuick 2.0
import Sailfish.Silica 1.0
import "../delegates"

Page {
    id: conversation
    objectName: "conversation"
    property bool editorFocus

    MessagesView {
        id: messages
        focus: true
        anchors.fill: parent

        model: MessageListModel

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
                    Backend.sendMessage(MessageModel.peerTel, text, "", attachmentPath)
                }
            }
        }
    }
}
