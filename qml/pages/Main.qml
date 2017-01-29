import QtQuick 2.0
import Sailfish.Silica 1.0
import "../delegates"

Page {
    id: main
    objectName: "main"

    Connections {
        target: Prompt
        onPromptPhoneNumber: {
            phoneNumberTimer.start()
        }
        onPromptVerificationCode: {
            verifyTimer.start()
        }
        onPromptPassword: {
            passwordTimer.start()
        }
        onPromptResetPeerIdentity: {
            pageStack.push(Qt.resolvedUrl("ResetPeerIdentity.qml"), { source: source })
        }
    }

    Connections {
        target: Backend
        onRegistrationSuccess: {
            registeredRemorse.execute("Registration complete!", function() { console.log("Registration complete") })
        }
    }

    RemorsePopup { id: registeredRemorse }

    SilicaListView {
        id: sessionView
        model: SessionListModel
        anchors.fill: parent
        spacing: Theme.paddingMedium

        PullDownMenu {
            MenuItem {
                text: qsTr("About Whisperfish")
                onClicked: pageStack.push(Qt.resolvedUrl("About.qml"))
            }
            MenuItem {
                text: qsTr("Settings")
                enabled: !Backend.locked
                onClicked: pageStack.push(Qt.resolvedUrl("Settings.qml"))
            }
            MenuItem {
                text: qsTr("New Message")
                enabled: !Backend.locked
                onClicked: pageStack.push(Qt.resolvedUrl("NewMessage.qml"))
            }
        }

        VerticalScrollDecorator {}

        ViewPlaceholder {
            enabled: sessionView.count == 0
            text: Backend.locked ? qsTr("Whisperfish") : qsTr("No messages")
            hintText: {
                if(!Backend.registered) {
                    qsTr("Registration required")
                } else if(Backend.locked) {
                    qsTr("Locked")
                } else {
                    ""
                }
            }
        }

        section {
            property: 'display.section'

            delegate: SectionHeader {
                text: section
                height: Theme.itemSizeExtraSmall
            }
        }

        delegate: Session{
            onClicked: {
                console.log("Activating session: "+model.display.id)
                pageStack.push(Qt.resolvedUrl("Conversation.qml"));
                SessionModel.markRead(model.display.id)
                MessageModel.refresh(
                    model.display.id,
                    Backend.contactName(model.display.source),
                    Backend.contactIdentity(model.display.source),
                    model.display.source,
                    model.display.isGroup
                )
            }
        }
    }

    Timer {
        id: phoneNumberTimer
        interval: 500
        running: false
        repeat: true
        onTriggered: {
            console.log("Page status: "+main.status)
            if(main.status == PageStatus.Active) {
                pageStack.push(Qt.resolvedUrl("Register.qml"))
                phoneNumberTimer.stop()
            }
        }
    }

    Timer {
        id: verifyTimer
        interval: 500
        running: false
        repeat: true
        onTriggered: {
            console.log("Page status: "+main.status)
            if(main.status == PageStatus.Active) {
                pageStack.push(Qt.resolvedUrl("Verify.qml"))
                verifyTimer.stop()
            }
        }
    }

    Timer {
        id: passwordTimer
        interval: 500
        running: false
        repeat: true
        onTriggered: {
            console.log("Page status: "+main.status)
            if(main.status == PageStatus.Active) {
                pageStack.push(Qt.resolvedUrl("Password.qml"))
                passwordTimer.stop()
            }
        }
    }
}
