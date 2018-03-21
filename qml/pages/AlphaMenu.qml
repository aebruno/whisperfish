/*
 * Author: r0kk3rz
 * https://github.com/r0kk3rz/sailfish-alphabet-sort
 */

import QtQuick 2.2
import Sailfish.Silica 1.0

Column {

    property Component listDelegate: listDelegate
    property ListModel dataSource: data

    width: parent.width
    id: contentColumn

        AlphaMenuGroupView {
            id: alphaGroupView
            width: parent.width
            opacity: enabled ? 1 : 0
            dataModel: dataSource
            Behavior on opacity { FadeAnimation {} }

            delegate: listDelegate

            onActivated: {
                // If height is reduced, allow the exterior flickable to reposition itself
                if (newViewHeight > alphaGroupView.height) {
                    // Where should the list be positioned to show as much of the list as possible
                    // (but also show one row beyond the list if possible)
                    var maxVisiblePosition = alphaGroupView.y + viewSectionY + newListHeight + rowHeight - parent.height

                    // Ensure up to two rows of group elements to show at the top
                    var maxAllowedPosition = alphaGroupView.y + Math.max(viewSectionY - (2 * rowHeight), 0)

                    // Don't position beyond the end of the flickable
                    var totalContentHeight = contentColumn.height + (newViewHeight - alphaGroupView.height)
                    var maxContentY = contentColumn.parent.originY + totalContentHeight - parent.height

                    var newContentY = Math.max(Math.min(Math.min(maxVisiblePosition, maxAllowedPosition), maxContentY), 0)
                    if (newContentY > parent.contentY) {
                        if (contentColumn.parent._contentYBeforeGroupOpen < 0) {
                            parent._contentYBeforeGroupOpen = parent.contentY
                        }
                        contentColumn.parent._animateContentY(newContentY, heightAnimationDuration, heightAnimationEasing)
                    }
                    console.log(alphaGroupView.height)
                }
                console.log(contentColumn.height)
            }
            onDeactivated: {
                if (contentColumn.parent._contentYBeforeGroupOpen >= 0) {
                    contentColumn.parent._animateContentY(root._contentYBeforeGroupOpen, heightAnimationDuration, heightAnimationEasing)
                    contentColumn.parent._contentYBeforeGroupOpen = -1
                }
            }
        }
}


