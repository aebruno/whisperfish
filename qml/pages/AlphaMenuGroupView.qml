/*
 * Author: r0kk3rz
 * https://github.com/r0kk3rz/sailfish-alphabet-sort
 */

import QtQuick 2.2
import Sailfish.Silica 1.0

Grid {
    id: root

    property int requiredProperty
    property Component delegate
    property int heightAnimationEasing: Easing.InOutQuad

    property real rowWidth: width / columns
    property real rowHeight: rowWidth
    property int rows: Math.ceil((groupModel.count + 1) / columns)

    property Item _currActiveGroup
    property Item _currentResultsList
    property Item _alternateResultsList

    property ListModel dataModel

    signal activated(real viewSectionY, real newListHeight, real newViewHeight, real heightAnimationDuration)
    signal deactivated(real heightAnimationDuration)

    columns: Math.floor(width / Theme.itemSizeMedium)

    function _groupAtIndex(index) {
        return groupsRepeater.itemAt(index)
    }

    function _lastIndexInRow(index) {
        var lastIndexInRow = ((Math.floor(
                                   index / columns) + 1) * root.columns) - 1
        var maxIndex = groupsRepeater.count
        return lastIndexInRow > maxIndex ? maxIndex : lastIndexInRow
    }

    function _getNextResultsList(parent) {
        if (_currentResultsList === null) {
            _currentResultsList = groupResultsListComponent.createObject(parent)
        } else if (_currentResultsList.parent === parent) {
            return _currentResultsList
        }

        if (_currentResultsList.active
                || (_alternateResultsList !== null
                    && _alternateResultsList.parent === parent)) {
            if (_alternateResultsList === null) {
                _alternateResultsList = _currentResultsList
                _currentResultsList = groupResultsListComponent.createObject(
                            parent)
            } else {
                var item = _alternateResultsList
                _alternateResultsList = _currentResultsList
                _currentResultsList = item
            }
        }

        if (_currentResultsList.parent !== null
                && _currentResultsList.parent !== parent)
            _currentResultsList.parent.groupResultsList = null
        _currentResultsList.parent = parent
        parent.groupResultsList = _currentResultsList
        return _currentResultsList
    }

    function _listForIndex(index) {
        var lastItemInRow = _groupAtIndex(_lastIndexInRow(index))
        if (_currentResultsList.parent === lastItemInRow)
            return _currentResultsList
        else if (_alternateResultsList.parent === lastItemInRow)
            return _alternateResultsList
        return null
    }

    function _openGroupList(name, index) {
        var list = _getNextResultsList(_groupAtIndex(_lastIndexInRow(index)))
        list.groupIndex = index

        list.open(name, index)
        return list.heightAnimationDuration
    }

    function _closeGroupList(index) {
        var list = _listForIndex(index)
        if (list) {
            list.close()
        }
    }

    function _groupListHeight(index) {
        var list = _listForIndex(index)
        if (list) {
            return list.implicitHeight
        }
        return 0
    }

    function _groupListOpenAnimationDuration(index, listCount) {
        var maxListItemsOnScreen = screen.height / Theme.itemSizeSmall
        if (listCount < maxListItemsOnScreen) {
            var minDuration = 150
            var maxDuration = 250
            return minDuration + ((maxDuration - minDuration) * (listCount / maxListItemsOnScreen))
        }
        // use default animation duration
        return 250
    }

    function _activate(group) {
        if ((_currentResultsList && _currentResultsList.animating)
                || (_alternateResultsList && _alternateResultsList.animating)) {
            // Wait til the previous animation completes before activating another
            return
        }
        if (group.active) {
            var list = _listForIndex(group.groupIndex)
            if (list) {
                deactivated(list.heightAnimationDuration)
            }
            _currActiveGroup = null
        } else if (group.hasEntries) {

            var heightAnimationDuration = _openGroupList(group.name,
                                                         group.groupIndex)
            if (_alternateResultsList !== null
                    && _alternateResultsList.active) {
                // the currently open list must close at the same rate as the new open list
                _alternateResultsList.heightAnimationDuration = heightAnimationDuration
            }
            _currActiveGroup = group

            var listHeight = _groupListHeight(group.groupIndex)
            activated((Math.floor(
                           group.groupIndex / columns) + 1) * group.baseHeight,
                      listHeight, (group.baseHeight * rows) + listHeight,
                      heightAnimationDuration)
        }
    }

    function _deactivate(group) {
        if (!group.active) {
            if ((_currActiveGroup == null) || Math.floor(
                        _currActiveGroup.groupIndex / columns) !== Math.floor(
                        group.groupIndex / columns)) {
                _closeGroupList(group.groupIndex)
            }
        }
    }

    onColumnsChanged: {
        if (_currActiveGroup !== null) {
            var oldLastItem = _currentResultsList.parent
            var newLastItem = _groupAtIndex(_lastIndexInRow(
                                                _currActiveGroup.groupIndex))
            if (oldLastItem !== newLastItem) {
                _currentResultsList.parent = newLastItem
                newLastItem.groupResultsList = _currentResultsList
                if (oldLastItem !== null)
                    oldLastItem.groupResultsList = null
            }
        }
    }

    //Empty List
    ListModel {
        id: emptyModel
    }

    //List to iterate over and build out top menu
    ListModel {
        id: groupModel

        ListElement { name: "A"; entryCount: 0 }
        ListElement { name: "B"; entryCount: 0 }
        ListElement { name: "C"; entryCount: 0 }
        ListElement { name: "D"; entryCount: 0 }
        ListElement { name: "E"; entryCount: 0 }
        ListElement { name: "F"; entryCount: 0 }
        ListElement { name: "G"; entryCount: 0 }
        ListElement { name: "H"; entryCount: 0 }
        ListElement { name: "I"; entryCount: 0 }
        ListElement { name: "J"; entryCount: 0 }
        ListElement { name: "K"; entryCount: 0 }
        ListElement { name: "L"; entryCount: 0 }
        ListElement { name: "M"; entryCount: 0 }
        ListElement { name: "N"; entryCount: 0 }
        ListElement { name: "O"; entryCount: 0 }
        ListElement { name: "P"; entryCount: 0 }
        ListElement { name: "Q"; entryCount: 0 }
        ListElement { name: "R"; entryCount: 0 }
        ListElement { name: "S"; entryCount: 0 }
        ListElement { name: "T"; entryCount: 0 }
        ListElement { name: "U"; entryCount: 0 }
        ListElement { name: "V"; entryCount: 0 }
        ListElement { name: "W"; entryCount: 0 }
        ListElement { name: "X"; entryCount: 0 }
        ListElement { name: "Y"; entryCount: 0 }
        ListElement { name: "Z"; entryCount: 0 }
        ListElement { name: "#"; entryCount: 0 }

        function countItems()
        {
            for(var i=0; (dataModel.count - 1) >= i; i++)
            {
                var index = dataModel.get(i).name.charAt(0)

                for(var j=0; (groupModel.count - 1 ) >= j; j++)
                {
                    if(groupModel.get(j).name === index.toUpperCase())
                    {
                        groupModel.setProperty(j, "entryCount", groupModel.get(j).entryCount + 1 )

                    }
                }
            }
        }

        Component.onCompleted: countItems();
    }


    Repeater {
        id: groupsRepeater
        model: groupModel

        ContactNameGroup {
            id: groupDelegate

            width: root.rowWidth
            baseHeight: root.rowHeight
            height: baseHeight + (groupResultsList !== null ? groupResultsList.height : 0)
            active: root._currActiveGroup === groupDelegate && hasEntries

            name: model.name
            groupIndex: model.index
            hasEntries: model.entryCount > 0

            onClicked: _activate(groupDelegate)
            onActiveChanged: _deactivate(groupDelegate)
        }
    }

    Component {
        id: groupResultsListComponent

        ColumnView {
            id: resultsView

            property real groupIndex
            property real heightAnimationDuration
            property bool animating: heightAnimation.running
                                     || fadeInAnimation.running

            property bool active: height > 0
            onActiveChanged: {
                if (!active) {
                    model = emptyModel
                }
            }

            function open(name, index) {
                filterModel.filterPattern = ''
                filterModel.filterPattern = name
                filterModel.filter()
                model = filterModel

                heightAnimationDuration = root._groupListOpenAnimationDuration(
                            groupIndex, model.count)
                if (state == "active") {
                    // already active, re-fade in with the new list contents
                    fadeInAnimation.start()
                }
                state = "active"
            }

            function close() {
                state = ""
            }

            itemHeight: Theme.itemSizeSmall

            model: emptyModel
            delegate: root.delegate
            cacheBuffer: itemHeight * 10


            //iterates through dataModel and checks name field against filterPattern
            //on match it adds item to filtermodel for display
            ListModel {
                id: filterModel
                property string filterPattern

                function filter() {
                    filterModel.clear()

                    for(var i=0; (dataModel.count - 1) >= i; i++)
                    {
                        filterPattern.charAt(0)
                        if(dataModel.get(i).name.charAt(0) === filterPattern.charAt(0) )
                        {
                            filterModel.append(dataModel.get(i))
                        }
                    }
                }
            }

            width: root.width
            height: 0
            x: -parent.x
            y: parent.baseHeight

            states: State {
                name: "active"
                PropertyChanges {
                    target: resultsView
                    height: resultsView.implicitHeight
                }
            }

            // use this instead of a Transition to animate the height because we need to trigger
            // this if open() is called when already active, and animating the height in a
            // standalone animation like fadeInAnimation will reset the height binding
            Behavior on height {
                enabled: !resultsView.menuOpen
                NumberAnimation {
                    id: heightAnimation
                    duration: resultsView.heightAnimationDuration
                    easing.type: root.heightAnimationEasing
                }
            }

            NumberAnimation {
                id: fadeInAnimation
                target: resultsView
                property: "opacity"
                from: 0.3
                to: 1
                duration: 300
                easing.type: Easing.InOutQuad
            }

            Rectangle {
                anchors.fill: parent
                z: parent.z - 1
                color: Theme.highlightBackgroundColor
                opacity: 0.1
            }
        }
    }
}
