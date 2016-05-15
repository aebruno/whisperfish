import QtQuick 2.0
import Sailfish.Silica 1.0

CoverBackground {
    Image {
        x: Theme.paddingLarge
        horizontalAlignment: Text.AlignHCenter
        source: sessionModel.unread > 0 ? "/usr/share/harbour-whisperfish/icons/86x86/harbour-whisperfish-gold.png" : "/usr/share/icons/hicolor/86x86/apps/harbour-whisperfish.png"
        anchors {
            bottom: parent.bottom
            bottomMargin: Theme.itemSizeLarge
            horizontalCenter: parent.horizontalCenter
        }
    }

    CoverActionList {
        CoverAction {
            iconSource: "image://theme/icon-cover-subview"
            onTriggered: {
                mainWindow.activate()
                showMainPage(PageStackAction.Immediate)
            }
        }

        CoverAction {
            iconSource: "image://theme/icon-cover-message"
            onTriggered: {
                if(!whisperfish.locked) {
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
