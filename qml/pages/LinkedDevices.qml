import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
    id: linkedDevices

    SilicaListView {
        id: listView
        anchors.fill: parent
        spacing: Theme.paddingMedium
        model: DeviceModel

        PullDownMenu {
            MenuItem {
                text: qsTr("Add")
                onClicked: {
                    var d = pageStack.push(Qt.resolvedUrl("AddDevice.qml"))
                    d.addDevice.connect(function(tsurl) {
                        console.log("Add device: "+tsurl)
                        // TODO: handle errors
                        DeviceModel.link(tsurl)
                    })
                }
            }
            MenuItem {
                text: qsTr("Refresh")
                onClicked: {
                    DeviceModel.reload()
                }
            }
        }
        header: PageHeader {
            title: qsTr("Linked Devices")
        }
        delegate: ListItem {
            contentHeight: created.y + created.height + lastSeen.height + Theme.paddingMedium
            id: delegate
            menu: deviceContextMenu

            function remove(contentItem) {
                contentItem.remorseAction(qsTr("Deleting"),
                    function() {
                        console.log("Unlink device: "+model.index)
                        DeviceModel.unlink(model.index)
                    })
            }

            Label {
                id: name
                truncationMode: TruncationMode.Fade
                font.pixelSize: Theme.fontSizeMedium
                text: model.name ? model.name : qsTr("Device "+model.id)
                anchors {
                    left: parent.left
                    leftMargin: Theme.horizontalPageMargin
                }
            }
            Label {
                function createdTime() {
                    var dt = new Date(model.created)
                    var linkDate = Format.formatDate(dt, Formatter.Timepoint)
                    return qsTr("Linked: "+linkDate)
                }
                id: created
                text: createdTime()
                font.pixelSize: Theme.fontSizeExtraSmall
                anchors {
                    top: name.bottom
                    left: parent.left
                    leftMargin: Theme.horizontalPageMargin
                }
            }
            Label {
                id: lastSeen
                function lastSeenTime() {
                    var dt = new Date(model.lastSeen)
                    var ls = Format.formatDate(dt, Formatter.DurationElapsed)
                    return qsTr("Last active: "+ls)
                }
                text: lastSeenTime()
                font.pixelSize: Theme.fontSizeExtraSmall
                font.italic: true
                anchors {
                    top: created.bottom
                    topMargin: Theme.paddingSmall
                    left: parent.left
                    leftMargin: Theme.horizontalPageMargin
                }
            }
            Component {
                id: deviceContextMenu
                ContextMenu {
                    id: menu
                    width: parent ? parent.width : Screen.width
                    MenuItem {
                        text: "Unlink"
                        onClicked: remove(menu.parent)
                    }
                }
            }
        }
    }
}
