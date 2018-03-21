import QtQuick 2.2
import Sailfish.Silica 1.0

CoverBackground {
    Image {
        x: Theme.paddingLarge
        horizontalAlignment: Text.AlignHCenter
        source: {
            if(SessionModel.unread > 0) {
                return "/usr/share/harbour-whisperfish/icons/86x86/harbour-whisperfish-gold.png"
            } else if(ClientWorker.connected) {
                return "/usr/share/harbour-whisperfish/icons/86x86/harbour-whisperfish-connected.png"
            } else if(!ClientWorker.connected) {
                return "/usr/share/harbour-whisperfish/icons/86x86/harbour-whisperfish-disconnected.png"
            } else {
                return "/usr/share/icons/hicolor/86x86/apps/harbour-whisperfish.png"
            }
        }
        anchors {
            bottom: parent.bottom
            bottomMargin: Theme.itemSizeHuge
            horizontalCenter: parent.horizontalCenter
        }
    }

    CoverActionList {
        CoverAction {
            iconSource: "image://theme/icon-cover-message"
            onTriggered: {
                if(!SetupWorker.locked) {
                    mainWindow.activate()
                    mainWindow.newMessage(PageStackAction.Immediate)
                }
            }
        }
    }

    Column {
        x: Theme.paddingLarge
        spacing: Theme.paddingSmall
        width: parent.width - 2*Theme.paddingLarge
        UnreadLabel {
            id: unreadLabel
        }
    }
}
