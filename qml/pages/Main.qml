import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
    id: main
    objectName: "main"

    property int sesslen: sessionModel.length

    onSesslenChanged: {
        refreshSessions()
    }

    property QtObject currentPage: pageStack.currentPage

    // This is a hack to use a psuedo model so we can use the 
    // group the messages into sections based on their timestamps
    function refreshSessions() {
        listView.model.clear()
        for (var i = 0; i < sessionModel.length; i++) {
            listView.model.append(sessionModel.get(i))
        }
    }

    function getPhoneNumber() {
        pageStack.push(Qt.resolvedUrl("Register.qml"))
    }

    function getVerificationCode() {
        pageStack.push(Qt.resolvedUrl("Verify.qml"))
    }

    function getStoragePassword() {
        pageStack.push(Qt.resolvedUrl("Password.qml"))
    }

    function registered() {
        registeredRemorse.execute("Registration complete!", function() { console.log("Registration complete") })
    }

    RemorsePopup { id: registeredRemorse }

    SilicaListView {
        id: listView
        model: ListModel {}
        anchors.fill: parent
        spacing: Theme.paddingMedium

        PullDownMenu {
            MenuItem {
                text: qsTr("About Whisperfish")
                onClicked: pageStack.push(Qt.resolvedUrl("About.qml"))
            }
            MenuItem {
                text: qsTr("New Message")
                onClicked: pageStack.push(Qt.resolvedUrl("NewMessage.qml"))
            }
        }

        VerticalScrollDecorator {}

        ViewPlaceholder {
            enabled: listView.count == 0
            text: "No messages"
            hintText: ""
        }

        section {
            property: 'section'

            delegate: SectionHeader {
                text: section
                height: Theme.itemSizeExtraSmall
            }
        }

        delegate: BackgroundItem {
            id: listItem
            width: parent.width
            height: Theme.itemSizeLarge

            property QtObject sess: sessionModel.get(index)

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
                    if(sent) {
                        "/usr/share/harbour-whisperfish/icons/ic_done_white_18dp.png"
                    } else if(recieved) {
                        "/usr/share/harbour-whisperfish/icons/ic_done_all_white_18dp.png"
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
                color: Theme.highlightColor
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
        }

        Component.onCompleted: {
            refreshSessions()
        }
    }
}
