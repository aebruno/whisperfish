import QtQuick 2.0
import Sailfish.Silica 1.0
import MeeGo.Connman 0.2
import org.nemomobile.contacts 1.0
import org.nemomobile.notifications 1.0
import "pages"

ApplicationWindow
{
    id: mainWindow
    cover: Qt.resolvedUrl("cover/CoverPage.qml")
    initialPage: Component { Main { } }
    allowedOrientations: Orientation.All
    _defaultPageOrientations: Orientation.All
    _defaultLabelFormat: Text.PlainText

    property bool connected: false
    property var enumFlags: {'group': 1, 'received': 2, 'unread': 4, 'sent': 8, 'attachment': 16, 'outgoing': 32}

    ImagePicker {
        id: imagepicker
    }

    PeopleModel {
        id: peopleModel
        filterType: PeopleModel.FilterNone
    }

    Component {
        id: messageNotification
        Notification {}
    }

    function activateSession(id) {
        var s = SessionModel.get(id)
        if(s) {
            SessionModel.markRead(s.id)
            MessageModel.refresh(
                s.id,
                s.name,
                Backend.contactIdentity(s.source),
                s.source,
                s.isGroup
            )
        }
    }

    function newMessageNotification(id, source, message) {
        if(Qt.application.state == Qt.ApplicationActive &&
           (pageStack.currentPage.objectName == "main" ||
           (id == MessageModel.sessionId && pageStack.currentPage.objectName == "conversation"))) {
           return
        }

        var m = messageNotification.createObject(null)
        if(SettingsBridge.boolValue("show_notify_message")) {
            m.body = message
        } else {
            m.body = qsTr("New Message")
        }
        m.category = "harbour-whisperfish-message"
        m.previewSummary = source
        m.previewBody = m.body
        m.summary = source
        m.clicked.connect(function() {
            console.log("Activating session: "+id)
            mainWindow.activate()
            showMainPage()
            pageStack.push(Qt.resolvedUrl("pages/Conversation.qml"), {}, PageStackAction.Immediate)
            mainWindow.activateSession(id)
        })
        // This is needed to call default action
        m.remoteActions = [ {
            "name": "default",
            "displayName": "Show Conversation",
            "icon": "harbour-whisperfish",
            "service": "org.whisperfish.session",
            "path": "/message",
            "iface": "org.whisperfish.session",
            "method": "showConversation",
            "arguments": [ "id", id ]
        } ]
        m.publish()
    }

    Connections {
        target: Backend
        onNotifyMessage: {
            newMessageNotification(id, source, message)
        }
    }

    Connections {
        target: SessionModel
        onRefresh: {
            SessionModel.load()
        }
        onMarkSent: {
            SessionModel.mark(sid, true, false, false)
        }
        onMarkReceived: {
            SessionModel.mark(sid, false, true, false)
        }
        onMarkRead: {
            SessionModel.mark(sid, false, false, true)
        }
        onUpdate: {
            if(sess.id == MessageModel.sessionId && pageStack.currentPage.objectName == "conversation") {
                sess.unread = false
            }
            SessionModel.add(sess)
        }
    }

    Connections {
        target: MessageModel
        onRefresh: {
            MessageModel.load(sid, peerName, peerIdentity, peerTel, group)
        }
        onMarkSent: {
            MessageModel.mark(mid, true, false)
        }
        onMarkReceived: {
            MessageModel.mark(mid, false, true)
        }
        onUpdate: {
            MessageModel.add(msg)
        }
    }

    function showMainPage() {
        pageStack.clear()
        pageStack.push(Qt.resolvedUrl("pages/Main.qml"), {}, PageStackAction.Immediate)
    }

    function newMessage(operationType) {
        showMainPage()
        pageStack.push(Qt.resolvedUrl("pages/NewMessage.qml"), { }, operationType)
    }

    function isConnected() {
        if(wifi.available && wifi.connected) {
            return true
        }
        if(cellular.available && cellular.connected) {
            return true
        }
        if(ethernet.available && ethernet.connected) {
            return true
        }

        return false
    }

    TechnologyModel {
        id: wifi
        name: "wifi"
        onConnectedChanged: {
            Backend.connected = mainWindow.isConnected()
        }
    }

    TechnologyModel {
        id: cellular
        name: "cellular"
        onConnectedChanged: {
            Backend.connected = mainWindow.isConnected()
        }
    }

    TechnologyModel {
        id: ethernet
        name: "ethernet"
        onConnectedChanged: {
            Backend.connected = mainWindow.isConnected()
        }
    }
}
