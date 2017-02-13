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
                //: Menu option to add new linked device
                //% "Add"
                text: qsTrId("whisperfish-add-linked-device")
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
                //: Menu option to refresh linked devices
                //% "Refresh"
                text: qsTrId("whisperfish-refresh-linked-devices")
                onClicked: {
                    DeviceModel.reload()
                }
            }
        }
        header: PageHeader {
            //: Title for Linked Devices page
            //% "Linked Devices"
            title: qsTrId("whisperfish-linked-devices")
        }
        delegate: ListItem {
            contentHeight: created.y + created.height + lastSeen.height + Theme.paddingMedium
            id: delegate
            menu: deviceContextMenu

            function remove(contentItem) {
                //: Unlinking remorse info message
                //% "Unlinking"
                contentItem.remorseAction(qsTrId("whisperfish-device-unlink-message"),
                    function() {
                        console.log("Unlink device: "+model.index)
                        DeviceModel.unlink(model.index)
                    })
            }

            Label {
                id: name
                truncationMode: TruncationMode.Fade
                font.pixelSize: Theme.fontSizeMedium
                text: model.name ? 
                    model.name : 
                //: Linked device name
                //% "Device %1"
                    qsTrId("whisperfish-device-name").arg(model.id)
                anchors {
                    left: parent.left
                    leftMargin: Theme.horizontalPageMargin
                }
            }
            Label {
                function createdTime() {
                    var dt = new Date(model.created)
                    var linkDate = Format.formatDate(dt, Formatter.Timepoint)
                    //: Linked device date
                    //% "Linked: %1"
                    return qsTrId("whisperfish-device-link-date").arg(linkDate)
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
                    //: Linked device last active date
                    //% "Last active: %1"
                    return qsTrId("whisperfish-device-last-active").arg(ls)
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
                        //: Device unlink menu option
                        //% "Unlink"
                        text: qsTrId("whisperfish-device-unlink")
                        onClicked: remove(menu.parent)
                    }
                }
            }
        }
    }
}
