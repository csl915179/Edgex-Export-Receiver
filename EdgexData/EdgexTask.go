package EdgexData

import contract "github.com/edgexfoundry/go-mod-core-contracts/models"

type EventList struct {
	List []contract.Event
}
func (L *EventList) Append(event contract.Event) {
	L.List = append(L.List,event)
}
func (L *EventList) Clear () []contract.Event {
	returnList := L.List[:]
	L.List = L.List[0:0]
	return returnList
}