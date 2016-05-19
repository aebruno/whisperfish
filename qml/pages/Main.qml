import QtQuick 2.0
import Sailfish.Silica 1.0
import "../delegates"

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
        var now = new Date().getTime()
        sessionView.model.clear()
        for (var i = 0; i < sessionModel.length; i++) {
            var s = sessionModel.get(i)
            var dt = new Date(s.timestamp)
            var sectionLabel = Format.formatDate(dt, Formatter.TimepointSectionRelative)
            sessionView.model.append({
                'id': s.id,
                'name': s.name,
                'isGroup': s.isGroup,
                'groupName': s.groupName,
                'received': s.received,
                'unread': s.unread,
                'sent': s.sent,
                'message': s.message,
                'hasAttachment': s.hasAttachment,
                'date': Format.formatDate(dt, Formatter.DurationElapsed),
                'section': sectionLabel ? sectionLabel : 'Today'
            })
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

    function confirmResetPeerIdentity(source) {
        pageStack.push(Qt.resolvedUrl("ResetPeerIdentity.qml"), { source: source })
    }

    RemorsePopup { id: registeredRemorse }

    SilicaListView {
        id: sessionView
        model: ListModel {}
        anchors.fill: parent
        spacing: Theme.paddingMedium

        PullDownMenu {
            MenuItem {
                text: qsTr("About Whisperfish")
                onClicked: pageStack.push(Qt.resolvedUrl("About.qml"))
            }
            MenuItem {
                text: qsTr("Settings")
                enabled: !whisperfish.locked
                onClicked: pageStack.push(Qt.resolvedUrl("Settings.qml"))
            }
            MenuItem {
                text: qsTr("New Message")
                enabled: !whisperfish.locked
                onClicked: pageStack.push(Qt.resolvedUrl("NewMessage.qml"))
            }
        }

        VerticalScrollDecorator {}

        ViewPlaceholder {
            enabled: sessionView.count == 0
            text: whisperfish.locked ? qsTr("Whisperfish") : qsTr("No messages")
            hintText: {
                if(!whisperfish.hasKeys) {
                    qsTr("Registration required")
                } else if(whisperfish.locked) {
                    qsTr("Locked")
                } else {
                    ""
                }
            }
        }

        section {
            property: 'section'

            delegate: SectionHeader {
                text: section
                height: Theme.itemSizeExtraSmall
            }
        }

        delegate: Session{}
    }
}
