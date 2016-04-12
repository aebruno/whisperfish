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
}
