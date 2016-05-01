import QtQuick 2.0
import Sailfish.Silica 1.0
import "pages"

ApplicationWindow
{
    id: mainWindow
    cover: Qt.resolvedUrl("cover/CoverPage.qml")
    initialPage: Component { Main { } }
    allowedOrientations: Orientation.All
    _defaultPageOrientations: Orientation.All
    _defaultLabelFormat: Text.PlainText

    ImagePicker {
        id: imagepicker
    }

    function showMainPage(operationType) {
        pageStack.clear()
        pageStack.push(Qt.resolvedUrl("pages/Main.qml"), {}, PageStackAction.Immediate)
    }

    function newMessage(operationType) {
        showMainPage(PageStackAction.Immediate)
        pageStack.push(Qt.resolvedUrl("pages/NewMessage.qml"), { }, operationType)
    }
}
