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
    }

    Connections {
        target: SetupWorker
        onRegistrationSuccess: {
            setupRemorse.execute("Registration complete!", function() { console.log("Registration complete") })
        }
        onInvalidDatastore: {
            setupRemorse.execute("Failed to setup datastore!", function() { console.log("Failed to setup datastore") })
        }
    }

    RemorsePopup { id: setupRemorse }

    SilicaListView {
        id: sessionView
        model: SessionModel
        anchors.fill: parent
        spacing: Theme.paddingMedium

        PullDownMenu {
            MenuItem {
                text: qsTr("About Whisperfish")
                onClicked: pageStack.push(Qt.resolvedUrl("About.qml"))
            }
            MenuItem {
                text: qsTr("Settings")
                enabled: !SetupWorker.locked
                onClicked: pageStack.push(Qt.resolvedUrl("Settings.qml"))
            }
            MenuItem {
                text: qsTr("New Message")
                enabled: !SetupWorker.locked
                onClicked: pageStack.push(Qt.resolvedUrl("NewMessage.qml"))
            }
        }

        VerticalScrollDecorator {}

        ViewPlaceholder {
            enabled: sessionView.count == 0
            text: SetupWorker.locked ? qsTr("Whisperfish") : qsTr("No messages")
            hintText: {
                if(!SetupWorker.registered) {
                    qsTr("Registration required")
                } else if(SetupWorker.locked) {
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

        delegate: Session{
            onClicked: {
                console.log("Activating session: "+model.id)
                pageStack.push(Qt.resolvedUrl("Conversation.qml"));
                if(model.unread) {
                    SessionModel.markRead(model.id)
                }
                MessageModel.load(
                    model.id,
                    ContactModel.name(model.source),
                    ContactModel.identity(model.source),
                    model.source,
                    model.isGroup
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
