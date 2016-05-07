import QtQuick 2.0
import Sailfish.Silica 1.0
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

    property var notifications: new Object()

    ImagePicker {
        id: imagepicker
    }

    Component {
        id: messageNotification
        Notification {}
    }

    function newMessageNotification(id, name, message) {
        var m = notifications[id]
        if(m) {
            m.itemCount++
        } else {
            m = messageNotification.createObject(null)
            m.itemCount = 1
        }
        var body = qsTr("New Message")
        if(whisperfish.settings().showNotifyMessage) {
            body = message
        } else if(m.itemCount > 1) {
            body += " ("+m.itemCount+")"
        }
        m.category = "harbour-whisperfish-message"
        m.previewSummary = name
        m.previewBody = body
        m.summary = name
        m.body = body
        m.clicked.connect(function() {
            mainWindow.activate()
            mainWindow.showSession(id, PageStackAction.Immediate)
        })
        // This is needed to call default action??
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
        notifications[id] = m
    }

    function showMainPage(operationType) {
        pageStack.clear()
        pageStack.push(Qt.resolvedUrl("pages/Main.qml"), {}, PageStackAction.Immediate)
    }

    function newMessage(operationType) {
        showMainPage(PageStackAction.Immediate)
        pageStack.push(Qt.resolvedUrl("pages/NewMessage.qml"), { }, operationType)
    }

    function removeNotification(id) {
        var n = notifications[id]
        if(n) {
            n.close()
            n.destroy()
            delete notifications[id]
        }
    }

    function showSession(id, operationType) {
        removeNotification(id)
        showMainPage(PageStackAction.Immediate)
        whisperfish.setSession(id)
        pageStack.push(Qt.resolvedUrl("pages/Conversation.qml"), { }, operationType)
    }
}
