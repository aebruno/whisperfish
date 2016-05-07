import QtQuick 2.0
import Sailfish.Silica 1.0

BackgroundItem {
    id: listItem
    width: parent.width
    height: Theme.itemSizeLarge

    Label {
        id: source
        text: isGroup ? qsTr('Group: '+groupName) : name
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
        mainWindow.removeNotification(id)
        whisperfish.setSession(id)
        pageStack.push(Qt.resolvedUrl("../pages/Conversation.qml"));
    }
}
