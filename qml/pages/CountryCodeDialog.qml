import QtQuick 2.0
import Sailfish.Silica 1.0
import "../js/iso_country_data.js" as CountryData

Dialog {
    canAccept: false
    signal setCountryCode(string text)

    SilicaListView {
        id: countryMenu
        anchors.fill: parent
        spacing: Theme.paddingMedium
        model: CountryData.isoCountries.length
        header: PageHeader {
            //: Directions for choosing country code
            //% "Choose Country Code"
            title: qsTrId("whisperfish-choose-country-code")
        }
        delegate: ListItem {
            Label {
                truncationMode: TruncationMode.Fade
                anchors {
                    left: parent.left
                    margins: Theme.paddingLarge
                }
                text: CountryData.isoCountries[index].ccode + " - " + CountryData.isoCountries[index].cname
            }
            onClicked: {
                setCountryCode(CountryData.isoCountries[index].ccode)
                pageStack.pop()
            }
        }
    }
}
