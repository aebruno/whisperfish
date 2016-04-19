import QtQuick 2.0
import Sailfish.Silica 1.0

BackgroundItem {
    id: listItem
    width: parent.width
    height: Theme.itemSizeLarge

    Label {
        id: source
        text: name
        font.pixelSize: Theme.fontSizeMedium
        truncationMode: TruncationMode.Fade
        anchors {
            left: parent.left
            right: status.left
            leftMargin: Theme.paddingLarge
        }
    }

    Image {
        source: {
            if(received) {
                "/usr/share/harbour-whisperfish/icons/ic_done_all_white_18dp.png"
            } else if(sent) {
                "/usr/share/harbour-whisperfish/icons/ic_done_white_18dp.png"
            } else {
                ""
            }
        }
        width: Theme.iconSizeSmall
        height: Theme.iconSizeSmall
        anchors {
            right: parent.right
            top: source.top
        }
    }

    Label {
        id: xbody
        text: message ? message : ''
        font.pixelSize: Theme.fontSizeExtraSmall
        wrapMode: Text.WordWrap
        maximumLineCount: 2
        color: unread ? Theme.highlightColor : Theme.primaryColor
        truncationMode: TruncationMode.Fade
        anchors {
            top: source.bottom
            left: parent.left
            right: parent.right
            leftMargin: Theme.paddingLarge
        }
    }

    Label {
        id: timestampLabel
        text: date
        font.pixelSize: Theme.fontSizeExtraSmall
        font.italic: true
        anchors {
            top: xbody.bottom
            topMargin: Theme.paddingSmall
            left: parent.left
            leftMargin: Theme.paddingLarge
            bottomMargin: Theme.paddingLarge
        }
    }


    onClicked: {
        whisperfish.setSession(id)
        sessionView.model.get(index).unread = false
        whisperfish.refreshConversation(id)
        messageModel.name = qsTr(""+name)
        messageModel.tel = qsTr(""+source)
        pageStack.push(Qt.resolvedUrl("../pages/Conversation.qml"));
    }
}
