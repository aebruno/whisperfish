import QtQuick 2.0
import Sailfish.Silica 1.0
import "../delegates"

Page {
    id: conversation
    objectName: "conversation"
    property bool editorFocus

    property int msglen: messageModel.length

    onMsglenChanged: {
        refreshMessages()
    }

    function updateSent(id) {
        for (var j = 0; j < messages.model.count; j++) {
            var m = messages.model.get(j)
            if(m.id == id) {
                messages.model.setProperty(j, "sent", true)
            }
        }
    }

    function updateReceived(id) {
        for (var j = 0; j < messages.model.count; j++) {
            var m = messages.model.get(j)
            if(m.id == id) {
                messages.model.setProperty(j, "received", true)
            }
        }
    }

    function refreshMessages() {
        var now = new Date().getTime()
        messages.model.clear()
        for (var i = 0; i < messageModel.length; i++) {
            var m = messageModel.get(i)
            var dt = new Date(m.timestamp)
            messages.model.append({
                'id': m.id,
                'sid': m.sid,
                'source': m.source,
                'message': m.message,
                'timestamp': m.timestamp,
                'outgoing': m.outgoing,
                'sent': m.sent,
                'received': m.received,
                'attachment': m.attachment,
                'mimeType': m.mimeType,
                'hasAttachment': m.hasAttachment
            })
        }
    }
    
    MessagesView {
        id: messages
        focus: true
        anchors.fill: parent

        model: ListModel {}

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
                contactName: messageModel.name
                enabled: true
                editorFocus: conversation.editorFocus

                onSendMessage: {
                    whisperfish.sendMessage(messageModel.tel, text, "", attachmentPath)
                }
            }
        }
    }
}
